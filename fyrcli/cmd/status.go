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

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the current status of the mesh.",
	Long: `Displays the current status of the mesh.

Prints out the values of the meshID (deviceID) and whether the mesh 
is currently set as being connected to the controller (meshconnected).`,

	Run: func(cmd *cobra.Command, args []string) {

		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("Connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Call the Status method.
		meshstatus, err := orch.Call_ORCH_Status(*client)
		if err != nil {
			fmt.Printf("Call to read mesh status failed -%v", err)
		}

		// Print the mesh status values.
		fmt.Printf("Mesh ID: %v\n", meshstatus.GetMeshID())
		fmt.Printf("Mesh Connected: %v\n", meshstatus.GetConnected())
	},
}

func init() {
	// Add the command 'status' to root CLI command.
	rootCmd.AddCommand(statusCmd)
}
