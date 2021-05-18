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

// nodelistCmd represents the nodelist command
var nodelistCmd = &cobra.Command{
	Use:   "nodelist",
	Short: "Displays the list of nodes connected to the mesh.",
	Long:  `Displays the list of nodes connected to the mesh.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("[error] connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Call the Status method.
		nodelist, err := orch.Call_ORCH_Nodelist(*client)
		if err != nil {
			fmt.Printf("[error] call to read mesh node list failed -%v", err)
		}

		// Iterate over the nodelist and print it.
		fmt.Println("mesh nodelist:")
		for index, node := range nodelist {
			fmt.Printf("%v] %v\n", index, node)
		}
	},
}

func init() {
	// Add the command 'nodelist' to root CLI command.
	rootCmd.AddCommand(nodelistCmd)
}
