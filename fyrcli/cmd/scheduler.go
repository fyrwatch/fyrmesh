/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	orch "github.com/fyrwatch/fyrmesh/fyrorch/orch"
	"github.com/spf13/cobra"
)

// schedulerCmd represents the scheduler command
var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Sets the state of the Scheduler",
	Long: `Sets the state of the Scheduler with a toggle.

The command requires the 'toggle' flag. The supported toggle phrases are given below:
- values such as 'on', 'true', 'play' and 'start' -> set the scheduler status to true.
- values such as 'off', 'false', 'pause' and 'stop' -> set the scheduler status to false. 

Setting the scheduler state will start/stop the scheduled ping activity.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the connection state set value from the command flags.
		toggle, _ := cmd.Flags().GetString("toggle")

		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("[error] connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Check the value of the connection state
		switch toggle {
		case "on", "true", "play", "start":
			// Call the SchedulerToggle method with the 'true' toggle.
			err := orch.Call_ORCH_SchedulerToggle(*client, true)
			// Check the error and print the appropriate message.
			if err == nil {
				fmt.Printf("[success] scheduler status successfully set to 'on'\n")
			} else {
				fmt.Printf("[failure] scheduler status failed to be set - %v\n", err)
			}

		case "off", "false", "pause", "stop":
			// Call the SchedulerToggle method with the 'true' toggle.
			err := orch.Call_ORCH_SchedulerToggle(*client, false)
			// Check the error and print the appropriate message.
			if err == nil {
				fmt.Printf("[success] scheduler status successfully set to 'off'\n")
			} else {
				fmt.Printf("[failure] scheduler status failed to be set - %v\n", err)
			}

		default:
			fmt.Println("[error] invalid value used for the 'toggle' flag!")
		}
	},
}

func init() {
	// Add the command 'scheduler' to root CLI command.
	rootCmd.AddCommand(schedulerCmd)

	// Add the flag 'toggle'
	schedulerCmd.Flags().StringP("toggle", "t", "", "value used to set the scheduler toggle.")
	schedulerCmd.MarkFlagRequired("toggle")
}
