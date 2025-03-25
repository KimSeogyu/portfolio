/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package internal

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "board",
	Short: "Board Server - A command line tool for managing board server",
	Long: `Board Server is a command line tool that helps you manage your board server.
It provides various commands to initialize, configure, and manage your board server.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
