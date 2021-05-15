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

	orch "github.com/fyrwatch/fyrmesh/fyrorch/orchpkg"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Set the connection state of the control node.",
	Long: `Set the connection state of the control node by sending the appropriate command to the it.

The command expects the 'set' flag which defaults to 'on'. The supported set phrases are given below:
- values such as 'on', 'true', 'connect' and 'yes' -> set the connection status to true.
- values such as 'off', 'false', 'disconnect' and 'no' -> set the connection status to false. 

Setting the connection state will simply flip the state of the onboard LED on the control node. 
This is useful to test whether the connection pipeline between the CLI and the control node is active, 
while also serving as a way to indicate that the CLI is currently communicating with the control node`,

	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the connection state set value from the command flags.
		set, _ := cmd.Flags().GetString("set")

		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("Connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Check the value of the connection state
		switch set {
		case "on", "true", "connect", "yes":
			// Call the Connection method with the 'true' connection state.
			success, err := orch.Call_ORCH_Connection(*client, true)
			// Check the acknowledgment and print the appropriate message.
			if success {
				fmt.Printf("Connection status successfully set to 'true'\n")
			} else {
				fmt.Printf("Connection status failed to be set - %v\n", err)
			}

		case "off", "false", "disconnect", "no":
			// Call the Connection method with the 'false' connection state.
			success, err := orch.Call_ORCH_Connection(*client, false)
			// Check the acknowledgment and print the appropriate message.
			if success {
				fmt.Printf("Connection status successfully set to 'off'\n")
			} else {
				fmt.Printf("Connection status failed to be set - %v\n", err)
			}

		default:
			fmt.Println("Invalid value used for the 'set' flag!")
		}
	},
}

func init() {
	// Add the command 'connect' to root CLI command.
	rootCmd.AddCommand(connectCmd)

	// Add the flag 'set'
	connectCmd.Flags().StringP("set", "s", "on", "value used to set the connection state.")
}
