/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"strings"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the daemon service",
	Long: `The 'stop' command is used to halt the execution of the CLI application. 
When this command is invoked, it prints "stop called" to the standard output and then terminates the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := exec.Command("pidof", "LiScreMon").CombinedOutput()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(strings.TrimSpace(string(result)))
		fmt.Println(os.Getpid())
		arrayOfpidByte := bytes.Split(result, []byte(" "))
		fmt.Println(arrayOfpidByte, len(arrayOfpidByte))

		if len(arrayOfpidByte) > 1 {
			pid := os.Getpid()
			thisProcess, err := strconv.Atoi(strings.TrimSpace(string(arrayOfpidByte[0])))
			if err != nil {
				log.Println("this should not happen", err)
			}
			if thisProcess == pid {
				for i := 1; i < len(arrayOfpidByte); i++ {
					otherProcess, err := strconv.Atoi(strings.TrimSpace(string(arrayOfpidByte[i])))
					if err != nil {
						log.Println("this should not happen2", err)
						continue
					}
					fmt.Println(otherProcess)
					err = syscall.Kill(otherProcess, syscall.SIGINT)
					if err != nil {
						log.Println("err killing", err)
					}

				}
				fmt.Println("success")
				return
			}

		}
		fmt.Println("LiScreMon was not running")

	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
