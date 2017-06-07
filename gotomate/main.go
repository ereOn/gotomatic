package main

import (
	"fmt"
	"io"
	"os"

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

		fmt.Println(config)

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&endpoint, "endpoint", "e", ":8080", "The endpoint to listen on")
	rootCmd.Flags().StringVarP(&configFile, "config-file", "c", "", "The configuration file to use")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
