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

// simulateCmd represents the simulate command
var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Starts a simulation of a Fire Event",
	Long:  `Starts a simulation of a Fire Event`,
	Run: func(cmd *cobra.Command, args []string) {

		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("[error] connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Call the SchedulerToggle method with the 'true' toggle.
		err = orch.Call_ORCH_Simulate(*client)
		// Check the error and print the appropriate message.
		if err == nil {
			fmt.Println("[success] fire event started successfully")
		} else {
			fmt.Println("[failure] fire event failed to start")
			fmt.Printf("[error] %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(simulateCmd)
}
