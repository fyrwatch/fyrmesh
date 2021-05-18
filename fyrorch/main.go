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
	logqueue := make(chan tools.Log)
	// Create an observer channel that will be used to pass observation logs.
	observerqueue := make(chan tools.ObserverLog)
	// Create a command queue that will be passed into the Orchestrator to siphon commands to the LINK.
	commandqueue := make(chan map[string]string)

	// Defer the closing of the created channels
	defer close(logqueue)
	defer close(observerqueue)
	defer close(commandqueue)

	//start a go-routine that handles log parsing, formatting, printing and forwarding.
	go tools.LogHandler(logqueue, observerqueue)

	// Initiate the connect runtime to the LINK server over gRPC
	client, conn, err := orch.GRPCconnect_LINK()
	defer conn.Close()
	if err != nil {
		logmessage := tools.NewOrchServerlog(fmt.Sprintf("connection to LINK server could not be established. error - %v", err))
		logqueue <- logmessage
	}

	// TODO: Setup Firebase Cloud Listener
	// TODO: Setup Task Generator and Scheduler

	// Start the go routine that initiates a read stream from the LINK server
	go orch.Call_LINK_Read(*client, logqueue)

	// Start the Orchestrator ORCH gRPC Server
	if err = orch.Start_ORCH_Server(*client, logqueue, commandqueue, observerqueue); err != nil {
		logmessage := tools.NewOrchServerlog(fmt.Sprintf("starting the ORCH server failed. error - %v", err))
		logqueue <- logmessage
	}
}
