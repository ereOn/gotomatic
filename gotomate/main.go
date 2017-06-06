package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	endpoint string
)

var rootCmd = &cobra.Command{
	Use:   "gotomate",
	Short: "Start an automation server.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.Flags().StringVarP(&endpoint, "endpoint", "e", ":8080", "The endpoint to listen on")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
