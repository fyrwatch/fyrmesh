/*
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.

This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed
without the express permission of Manish Meganathan and Mariyam A.Ghani
===========================================================================
FyrMesh FyrORCH
===========================================================================
*/
package main

import (
	"fmt"

	orch "github.com/fyrwatch/fyrmesh/fyrorch/orch"
	tools "github.com/fyrwatch/fyrmesh/tools"
)

func main() {
	// Construct a new MeshOrchestrator
	meshorchestrator, err := tools.NewMeshOrchestrator()
	if err != nil {
		// Generate an ORCH serverlog and print it.
		fmt.Println(tools.FormatLog(tools.NewOrchServerlog(fmt.Sprintf("(error) mesh orchestrator could not be constructed | error - %v |", err))))
	}

	// Defer the closing of the meshorchestrator channels
	defer meshorchestrator.Close()

	// Start a go-routine that handles log parsing, formatting, printing and forwarding.
	go tools.LogHandler(meshorchestrator)

	// Initiate the connect runtime to the LINK server over gRPC
	client, conn, err := orch.GRPCconnect_LINK()
	defer conn.Close()
	if err != nil {
		// Generate an ORCH serverlog and send it over the LogQueue of the meshorchestrator
		fmt.Println(tools.FormatLog(tools.NewOrchServerlog(fmt.Sprintf("(error) connection to LINK server could not be established | error - %v |", err))))
	}

	// Start the go routine that starts streaming logs from the LINK server
	go orch.Call_LINK_Read(*client, meshorchestrator.LogQueue)

	// Start the Orchestrator ORCH gRPC Server
	if err = orch.Start_ORCH_Server(*client, meshorchestrator); err != nil {
		// Generate an ORCH serverlog and send it over the LogQueue of the meshorchestrator
		fmt.Println(tools.FormatLog(tools.NewOrchServerlog(fmt.Sprintf("(error) starting the ORCH server failed | error - %v |", err))))
	}
}
