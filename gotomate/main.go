package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/intelux/gotomatic/conditional"
	"github.com/intelux/gotomatic/configuration"
	"github.com/spf13/cobra"
)

var (
	endpoint   string
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "gotomate",
	Short: "Start an automation server.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var config configuration.Configuration

		if configFile != "" {
			var f io.ReadCloser

			if f, err = os.Open(configFile); err != nil {
				return err
			}

			defer f.Close()

			if config, err = configuration.Load(f); err != nil {
				return err
			}
		} else {
			config = configuration.New()
		}

		defer config.Close()

		r := mux.NewRouter()
		r.Methods("GET").Path("/conditions/{name}").HandlerFunc(GetConditionHandler(config))
		r.Methods("POST").Path("/conditions/{name}").HandlerFunc(WaitConditionHandler(config))
		r.Methods("PUT").Path("/conditions/{name}").HandlerFunc(SetConditionHandler(config))

		stop := make(chan os.Signal)
		defer close(stop)

		signal.Notify(stop, os.Interrupt)

		server := &http.Server{Addr: endpoint, Handler: r}
		errorCh := make(chan error, 1)

		go func() {
			defer close(errorCh)

			if err := server.ListenAndServe(); err != nil {
				errorCh <- err
			}
		}()

		fmt.Printf("Started HTTP server on %s.\n", endpoint)
		defer fmt.Printf("Stopped HTTP server.\n")

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		go func() {
			if err := config.Watch(ctx); err != nil {
				errorCh <- err
			}
		}()

		select {
		case <-stop:
			fmt.Printf("Interruption. Exiting...\n")
			cancel()

			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()

			server.Shutdown(ctx)
			return server.Close()
		case err := <-errorCh:
			return err
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&endpoint, "endpoint", "e", ":8080", "The endpoint to listen on")
	rootCmd.Flags().StringVarP(&configFile, "config-file", "c", "", "The configuration file to use")
}

func GetConditionHandler(config configuration.Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name := mux.Vars(req)["name"]
		condition := config.GetCondition(name)

		if condition == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		state, _ := condition.GetAndWaitChange()
		json.NewEncoder(w).Encode(state)
	}
}

func WaitConditionHandler(config configuration.Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name := mux.Vars(req)["name"]
		condition := config.GetCondition(name)

		if condition == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var state bool

		if err := json.NewDecoder(req.Body).Decode(&state); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s\n", fmt.Errorf("JSON decoding error: %s", err))
			return
		}

		ch := condition.Wait(state)

		select {
		case <-req.Context().Done():
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusServiceUnavailable)
		case <-ch:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(state)
		}
	}
}

func SetConditionHandler(config configuration.Configuration) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name := mux.Vars(req)["name"]
		condition := config.GetCondition(name)

		if condition == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		settable, ok := condition.(conditional.Settable)

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "unsettable condition type\n")
			return
		}

		var state bool

		if err := json.NewDecoder(req.Body).Decode(&state); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s\n", fmt.Errorf("JSON decoding error: %s", err))
			return
		}

		settable.Set(state)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
