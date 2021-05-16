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
	// Create a log channel that will be used to pass all logs within the server.
	logqueue := make(chan string)
	// Create an observer channel that will be used to pass observation logs.
	obsqueue := make(chan string)
	// Create a command queue that will be passed into the Orchestrator to siphon commands to the LINK.
	commandqueue := make(chan string)

	// Defer the closing of the created channels
	defer close(logqueue)
	defer close(obsqueue)
	defer close(commandqueue)

	//start a go-routine that handles log printing and forwarding.
	go tools.LogHandler(logqueue, obsqueue)

	// Initiate the connect runtime to the LINK server over gRPC
	client, conn, err := orch.GRPCconnect_LINK()
	defer conn.Close()
	if err != nil {
		logmessage := tools.GenerateORCHLog(fmt.Sprintf("connection to interface LINK failed - %v", err))
		logqueue <- logmessage
	}

	// TODO: Setup Firebase Cloud Listener
	// TODO: Setup Task Generator and Scheduler

	// Start the go routine that initiates a read stream from the LINK server
	go orch.Call_LINK_Read(*client, logqueue)

	// Start the Orchestrator ORCH gRPC Server
	if err = orch.Start_ORCH_Server(*client, logqueue, commandqueue, obsqueue); err != nil {
		logmessage := tools.GenerateORCHLog(fmt.Sprintf("serving orchestrator ORCH failed - %v", err))
		logqueue <- logmessage
	}
}
