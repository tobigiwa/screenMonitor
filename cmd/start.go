/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"LiScreMon/daemon"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "launch the daemon service",
	Long: `Command "start" launch the daemon service. It starts the background screen
	monitoring and recording program of LiScreMon`,
	Run: func(cmd *cobra.Command, args []string) {
		daemon.DaemonService()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
