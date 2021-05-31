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
			fmt.Printf("[error] connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Call the Status method.
		meshstatus, err := orch.Call_ORCH_Status(*client)
		if err != nil {
			fmt.Printf("[error] call to read mesh status failed -%v", err)
		}

		// Retrieve the map of nodes
		nodelist := meshstatus.GetNodelist().GetNodes()

		// Print the mesh status values.
		fmt.Printf("mesh connection state: %v\n", meshstatus.GetConnected())
		fmt.Println()
		fmt.Printf("mesh controller ID: %v\n", meshstatus.GetControllerID())
		fmt.Printf("mesh controlnode ID: %v\n", meshstatus.GetControlnodeID())
		fmt.Println()
		fmt.Printf("mesh SSID: %v\n", meshstatus.GetMeshSSID())
		fmt.Printf("mesh PORT: %v\n", meshstatus.GetMeshPORT())
		fmt.Printf("mesh password: %v\n", meshstatus.GetMeshPSWD())
		fmt.Println()
		fmt.Println("mesh nodelist:")

		index := 1
		for nodeid, nodeconfig := range nodelist {
			fmt.Printf("%v] %v\t%v\n", index, nodeid, nodeconfig)
			index++
		}
	},
}

func init() {
	// Add the command 'status' to root CLI command.
	rootCmd.AddCommand(statusCmd)
}
