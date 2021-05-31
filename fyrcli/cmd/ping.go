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

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Pings the mesh or a node.",
	Long: `Pings the mesh or a node for sensor/config data by sending the appropriate control command.

The 'type(t)' flag sets the type of ping to perform. Valid values are 'sensor', 'config' and 'control'.
The 'node(n)' flag sets a node ID to ping. If this value is not, the whole mesh is pinged.
The 'phrase(p)' flag sets the user phrase to use in the ping ID, which has the format 'userping-<phrase>' for user generated pings.

If the type is set as 'control', the other flags are ignored and the control node is pinged for its config. 

Note: the responses from the pings are currently not captured and will appear in the 
ORCH logs with a MESH source and a sensordata/configdata type.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the command flags
		pingtype, _ := cmd.Flags().GetString("type")
		pingnode, _ := cmd.Flags().GetString("node")
		pingphrase, _ := cmd.Flags().GetString("phrase")
		// Declare a trigger string
		var pingtrigger string

		// Check the value of the pingtype
		switch pingtype {
		case "sensor", "config":
			if pingnode == "mesh" {
				pingtrigger = fmt.Sprintf("ping-%s-mesh", pingtype)
			} else {
				pingtrigger = fmt.Sprintf("ping-%s-node", pingtype)
			}

		case "control":
			pingtrigger = "ping-control"

		default:
			fmt.Println("[error] an invalid ping type was provided")
		}

		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("[error] connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Call the Ping method with the trigger, node and phrase.
		success, err := orch.Call_ORCH_Ping(*client, pingtrigger, pingnode, pingphrase)
		// Check the acknowledgment and print the appropriate message.
		if success {
			fmt.Printf("[success] %v was pinged successfully\n", pingnode)
		} else {
			fmt.Printf("[failure] %v was failed to be pinged\n", pingnode)
			fmt.Printf("[error] %v\n", err)
		}
	},
}

func init() {
	// Add the command 'ping' to root CLI command.
	rootCmd.AddCommand(pingCmd)

	// Add the flag 'type' and mark as required
	pingCmd.Flags().StringP("type", "t", "sensor", "type of ping to perform")
	// Add the flag 'node'
	pingCmd.Flags().StringP("node", "n", "mesh", "node ID to ping")
	// Add the flag 'phrase'
	pingCmd.Flags().StringP("phrase", "p", "atestping", "phrase to use with ping ID")
}
