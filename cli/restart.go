/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"LiScreMon/cli/daemon"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// backgroundCmd represents the background command
var backgroundCmd = &cobra.Command{
	Use:   "restart",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := stopLiscrenMon(); err != nil {
			log.Println("could not sucessfully shutdown running LiScreMon", err)
			return
		}
		fmt.Println("LiScreMon would be restarting now...")
		time.Sleep(1 * time.Second) // allow for all resources to be released
		daemon.DaemonServiceLinux()
	},
}

func init() {
	rootCmd.AddCommand(backgroundCmd)
}
