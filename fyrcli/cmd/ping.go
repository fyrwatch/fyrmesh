/*
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.

This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed
without the express permission of Manish Meganathan and Mariyam A.Ghani
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
