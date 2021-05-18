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

	"github.com/spf13/cobra"

	orch "github.com/fyrwatch/fyrmesh/fyrorch/orch"
)

// commandCmd represents the command command
var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "Sends a control command to the mesh.",
	Long: `Sends a control command to the mesh control node. 
The message flag is mandatory and is used to set the control command phrase.
All other arguments are collected as key value pairs for the command metadata. 
Metadata collection is only done if an even number of args are provided.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the message value from the command flags.
		message, _ := cmd.Flags().GetString("message")
		// Set the command message into a command map
		commandmap := map[string]string{"command": message}

		// Check if any args have been passed to convert into metadata
		if len(args) != 0 {
			if len(args)%2 != 0 {
				// An uneven number of args cannot form key value pairs
				fmt.Println("[error] an uneven number of metadata key-value pair arguments were provided")
				return
			}

			// Iterate over the arguments and collect them into the map
			for i := 0; i < len(args); i += 2 {
				commandmap[args[i]] = args[i+1]
			}
		}

		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("[error] connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Call the Command method with the commandmap
		success, err := orch.Call_ORCH_Command(*client, commandmap)
		// Check the success value and print output
		if success {
			fmt.Println("[success] command was sent successfully")
		} else {
			fmt.Println("[failure] command failed to be sent")
			fmt.Printf("[error] %v", err)
		}
	},
}

func init() {
	// Add the command 'command' to root CLI command.
	rootCmd.AddCommand(commandCmd)

	// Add the flag 'command' and mark as required.
	commandCmd.Flags().StringP("message", "m", "", "command message to send")
	commandCmd.MarkFlagRequired("message")

	// Define the usage template
	usage := `Usage:
fyrcli command -m [message] metadataKey1 metadataVal1...

Flags:
-h, --help             help for command
-m, --message string   command message to send
`
	// Set the custom usage template
	commandCmd.SetUsageTemplate(usage)
}
