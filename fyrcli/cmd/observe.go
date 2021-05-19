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

	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"

	orch "github.com/fyrwatch/fyrmesh/fyrorch/orch"
	tools "github.com/fyrwatch/fyrmesh/tools"
)

// observeCmd represents the observe command
var observeCmd = &cobra.Command{
	Use:   "observe",
	Short: "Observes the logstream from the ORCH server.",
	Long: `Observes the logstream from the ORCH server and prints them to the console.

The observation of the logstream can be filtered based on the source of the log, 
the type of the log or both. A valid combination must be used when both are passed.

Valid source and type filters:
- For the source filter:	'MESH', 'LINK' and 'ORCH'.
- For the type filter:
	- 'serverlog', 'protolog' (only supported for the 'LINK' and 'ORCH' source filter)
	- 'cloudlog', 'schedlog', 'obstoggle' (only supported for the 'ORCH' source filter)
	- 'message', 'newconnection', 'changedconnection', 'nodetimeadjust', 'handshake', 'sensordata',
	'configdata', 'controlconfig', 'nodelist' (only supported for the 'MESH' source filter)

Observation of ORCH log can only performed by a device configured as a 'mesh-observer'.
The observer collects the logs being printed to the ORCH server console and prints them on the 
terminal that invokes it. Observer logs have the '[OBS]' suffix followed by the log itself.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Read the config file.
		config, err := tools.ReadConfig()
		if err != nil {
			fmt.Printf("[error] config file could not be read - %v\n", err)
			fmt.Println("[suggestion] run 'fyrcli config -m generate' if file does not exist or is corrupted.")
			return
		}

		// Check the device type config value.
		if config.DeviceType != "mesh-observer" {
			fmt.Println("[error] log observation can only be performed by a mesh observer.")
			return
		}

		// Connect to the ORCH gRPC server.
		client, conn, err := orch.GRPCconnect_ORCH()
		defer conn.Close()
		if err != nil {
			fmt.Printf("[error] connection to ORCH gRPC server could not be established - %v\n", err)
		}

		// Retrieve the log filter flags
		sourcefilter, _ := cmd.Flags().GetString("source")
		typefilter, _ := cmd.Flags().GetString("type")

		// Check the value of source filter
		switch sourcefilter {
		case "ORCH":
			// Check the value of type filter
			switch typefilter {
			case "serverlog", "protolog", "cloudlog", "schedlog", "obstoggle", "":
			default:
				fmt.Println("[error] invalid type filter applied for the 'ORCH' source filter")
				return
			}

		case "LINK":
			// Check the value of type filter
			switch typefilter {
			case "serverlog", "protolog", "":
			default:
				fmt.Println("[error] invalid type filter applied for the 'LINK' source filter")
				return
			}

		case "MESH":
			// Check the value of type filter
			switch typefilter {
			case "message", "newconnection", "changedconnection", "nodetimeadjust", "":
			case "handshake", "sensordata", "configdata", "controlconfig", "nodelist":
			default:
				fmt.Println("[error] invalid type filter applied for the 'MESH' source filter")
				return
			}

		case "":
			switch typefilter {
			case "serverlog", "protolog", "cloudlog", "schedlog", "obstoggle", "":
			case "message", "newconnection", "changedconnection", "nodetimeadjust":
			case "handshake", "sensordata", "configdata", "controlconfig", "nodelist":
			default:
				fmt.Println("[error] invalid type filter applied")
				return
			}

		default:
			fmt.Println("[error] invalid source filter applied")
			return
		}

		// Call the Observe method of the ORCH server with the log filters
		stream, err := orch.Call_ORCH_Observe(*client, sourcefilter, typefilter)
		if err != nil {
			fmt.Printf("[error] observe stream failed to be established - %v\n", err)
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
				fmt.Printf("[error] observe stream broke. error while streaming - (%v)%v", errstatus.Code(), errstatus.Message())
				break
			}

			// Print the observer log to the console.
			fmt.Printf("[OBS] %v\n", observelog.GetMessage())
		}
	},
}

func init() {
	// Add the command 'observe' to root CLI command.
	rootCmd.AddCommand(observeCmd)

	// Add the flag 'source'
	observeCmd.Flags().StringP("source", "s", "", "value of source log filter")
	// Add the flag 'type'
	observeCmd.Flags().StringP("type", "t", "", "value of type log filter")
}
