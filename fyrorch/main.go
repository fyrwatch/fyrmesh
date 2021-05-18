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
		// Generate an ORCH serverlog and print it. (The logqueue was never built and hence cannot it be pushed into)
		logmessage := tools.NewOrchServerlog(fmt.Sprintf("mesh orchestrator could not be constructed. error - %v", err))
		fmt.Println(tools.StringifyLog(logmessage))
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
		logmessage := tools.NewOrchServerlog(fmt.Sprintf("connection to LINK server could not be established. error - %v", err))
		meshorchestrator.LogQueue <- logmessage
	}

	// Start the go routine that starts streaming logs from the LINK server
	go orch.Call_LINK_Read(*client, meshorchestrator.LogQueue)

	// TODO: Setup Firebase Cloud Listener
	// TODO: Setup Task Generator and Scheduler

	// Start the Orchestrator ORCH gRPC Server
	if err = orch.Start_ORCH_Server(*client, meshorchestrator); err != nil {
		// Generate an ORCH serverlog and send it over the LogQueue of the meshorchestrator
		logmessage := tools.NewOrchServerlog(fmt.Sprintf("starting the ORCH server failed. error - %v", err))
		meshorchestrator.LogQueue <- logmessage
	}
}
