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
	"io"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"

	orch "github.com/fyrwatch/fyrmesh/fyrorch/orch"
	tools "github.com/fyrwatch/fyrmesh/tools"
)

// observeCmd represents the observe command
var observeCmd = &cobra.Command{
	Use:   "observe",
	Short: "Observes the logstream of the ORCH server.",
	Long: `Observes the logstream of the ORCH server and prints them to the console.

Observation of ORCH log can only performed by a device configured as a 'mesh-observer'.
The observer collects the logs being printed to the ORCH server console and prints them on the 
terminal that invokes it. Observer logs have the '[OBS]' suffix followed by the log itself.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Read the config file.
		config, err := tools.ReadConfig()
		if err != nil {
			fmt.Printf("Config file could not be read - %v\n", err)
			fmt.Println("Run 'fyrcli config -m generate' if file does not exist or is corrupted.")
			os.Exit(0)
		}

		// Check the device type config value.
		if config.DeviceType != "mesh-observer" {
			fmt.Println("Log observation can only be performed by a mesh observer.")
			os.Exit(0)
		}

		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("Connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Call the Observe method of the ORCH server.
		stream, err := orch.Call_ORCH_Observe(*client)
		if err != nil {
			fmt.Printf("Observe stream failed to be established - %v\n", err)
		}

		// Start an infinite loop to read from the stream
		for {
			// Recieve an Message object from the stream
			observelog, err := stream.Recv()

			// Break out of loop if stream has closed
			if err == io.EOF {
				break
			}

			// Print any other error and break out of the loop.
			if err != nil {
				errstatus, _ := status.FromError(err)
				fmt.Printf("Observe stream broke. error while streaming - (%v)%v", errstatus.Code(), errstatus.Message())
				break
			}

			// Print the observer log to the console.
			fmt.Printf("[OBS]%v\n", observelog.GetMessage())
		}
	},
}

func init() {
	// Add the command 'observe' to root CLI command.
	rootCmd.AddCommand(observeCmd)
}
