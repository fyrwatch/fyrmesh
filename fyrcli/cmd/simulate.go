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
