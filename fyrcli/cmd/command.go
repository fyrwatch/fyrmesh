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
			fmt.Printf("[error] %v\n", err)
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
