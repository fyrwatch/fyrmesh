/*
===========================================================================
MIT License

Copyright (c) 2021 Manish Meganathan, Mariyam A.Ghani

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
===========================================================================
FyrMesh FyrCLI
===========================================================================
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
