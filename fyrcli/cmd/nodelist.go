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

		// Call the Nodelist method.
		nodelist, err := orch.Call_ORCH_Nodelist(*client)
		if err != nil {
			fmt.Printf("[error] call to read mesh node list failed -%v", err)
		}

		// Iterate over the nodelist and print it.
		fmt.Println("mesh nodelist:")

		index := 1
		for nodeid, nodeconfig := range nodelist {
			fmt.Printf("%v] %v\t%v\n", index, nodeid, nodeconfig)
			index++
		}
	},
}

// nodelistCmd represents the nodelist command
var nodelistUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the list of nodes connected to the mesh.",
	Long:  `Updates the list of nodes connected to the mesh.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("[error] connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Call the Command method.
		commandmap := map[string]string{"command": "readnodelist-control"}
		success, err := orch.Call_ORCH_Command(*client, commandmap)
		if err != nil {
			fmt.Printf("[error] call to update mesh node list failed -%v", err)
		}

		// Check the success value and print output
		if success {
			fmt.Println("[success] command to update nodelist was sent successfully")
		} else {
			fmt.Println("[failure] command to update nodelist failed to be sent")
			fmt.Printf("[error] %v\n", err)
		}
	},
}

func init() {
	// Add the command 'nodelist' to root CLI command.
	rootCmd.AddCommand(nodelistCmd)
	//
	nodelistCmd.AddCommand(nodelistUpdateCmd)
}
