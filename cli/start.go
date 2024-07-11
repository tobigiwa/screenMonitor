/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"LiScreMon/daemon"
	"log"
	"log/slog"
	helperFuncs "pkg/helper"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts the daemon service",
	Long: `Command "start" launch the daemon service. It starts the background screen
	monitoring and recording program of LiScreMon`,
	Run: func(cmd *cobra.Command, args []string) {
		// cpuProfileFile, err := os.Create("cpuProfile.prof")
		// if err != nil {
		// 	panic(err)
		// }
		// defer cpuProfileFile.Close()
		// if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
		// 	panic(err)
		// }
		// defer pprof.StopCPUProfile()

		// memoryProfileFile, err := os.Create("memoryProfile.prof")
		// if err != nil {
		// 	panic(err)
		// }
		// defer memoryProfileFile.Close()
		// if err := pprof.StartCPUProfile(memoryProfileFile); err != nil {
		// 	panic(err)
		// }
		// defer pprof.StopCPUProfile()
		logger, logFile, err := helperFuncs.Logger("daemon.log")
		if err != nil {
			log.Fatalln(err) // exit
		}
		defer logFile.Close()

		slog.SetDefault(logger)

		daemon.DaemonServiceLinux(logger)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
