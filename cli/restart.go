/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"
	"log"
	"os/exec"

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
		err := exec.Command("LiScreMon", "start").Start()
		if err != nil {
			log.Println("error starting", err)
			return
		}
		fmt.Println("LiSreMon is running again")
	},
}

func init() {
	rootCmd.AddCommand(backgroundCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backgroundCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backgroundCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
