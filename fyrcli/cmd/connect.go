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
			fmt.Printf("[error] connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Check the value of the connection state
		switch set {
		case "on", "true", "connect", "yes":
			// Call the Connection method with the 'true' connection state.
			success, err := orch.Call_ORCH_Connection(*client, true)
			// Check the acknowledgment and print the appropriate message.
			if success {
				fmt.Printf("[success] connection status successfully set to 'true'\n")
			} else {
				fmt.Printf("[failure] connection status failed to be set - %v\n", err)
			}

		case "off", "false", "disconnect", "no":
			// Call the Connection method with the 'false' connection state.
			success, err := orch.Call_ORCH_Connection(*client, false)
			// Check the acknowledgment and print the appropriate message.
			if success {
				fmt.Printf("[success] connection status successfully set to 'off'\n")
			} else {
				fmt.Printf("[failure] connection status failed to be set - %v\n", err)
			}

		default:
			fmt.Println("[error] invalid value used for the 'set' flag!")
		}
	},
}

func init() {
	// Add the command 'connect' to root CLI command.
	rootCmd.AddCommand(connectCmd)

	// Add the flag 'set'
	connectCmd.Flags().StringP("set", "s", "on", "value used to set the connection state.")
}
