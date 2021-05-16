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

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Pings the mesh for sensor readings.",
	Long: `Pings the mesh for sensor readings by sending a command to the mesh controller.

The responses from the pings are currently not captured and will appear in the 
ORCH logs like any other message from the mesh with the [MESH] suffix attached.`,

	Run: func(cmd *cobra.Command, args []string) {

		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("Connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Call the Ping method to ping the whole mesh.
		success, err := orch.Call_ORCH_Ping(*client)
		// Check the acknowledgment and print the appropriate message.
		if success {
			fmt.Printf("Mesh was pinged successfully\n")
		} else {
			fmt.Printf("Mesh was failed to be pinged - %v\n", err)
		}
	},
}

func init() {
	// Add the command 'ping' to root CLI command.
	rootCmd.AddCommand(pingCmd)

	// Add the flag 'node'
	pingCmd.Flags().StringP("node", "n", "", "node identifier to ping (deprecated)")
}
