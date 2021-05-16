/*
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.

This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed
without the express permission of Manish Meganathan and Mariyam A.Ghani
===========================================================================
FyrMesh FyrORCH gopkg orch
===========================================================================
*/
package orch

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/fyrwatch/fyrmesh/proto"
	tools "github.com/fyrwatch/fyrmesh/tools"
)

// A function that establishes a gRPC connection to the orchestrator ORCH server and returns the ORCH
// client object, the gPRC connection object and any error that occurs while attempting to connect.
func GRPCconnect_ORCH() (*pb.OrchestratorClient, *grpc.ClientConn, error) {
	// Read the service config for the ORCH server
	config, err := tools.ReadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("could not read config file - %v", err)
	}

	// Construct the URL for the orchestrator ORCH server
	orchconfig := config.Services["ORCH"]
	orchhost := fmt.Sprintf("%s:%d", orchconfig.Host, orchconfig.Port)

	// Connect to Orchestrator ORCH gRPC Server
	conn, err := grpc.Dial(orchhost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("could not dialup ORCH gRPC server - %v", err)
	}

	// Create an Orchestrator ORCH Client and return it along with the connection object and a nil error
	client := pb.NewOrchestratorClient(conn)
	return &client, conn, nil
}

// A function that calls the 'Connection' method of the ORCH server over a gRPC connection.
// Requires the ORCH client object and a boolean value of the connection state to transmit.
func Call_ORCH_Connection(client pb.OrchestratorClient, value bool) (bool, error) {
	// Test the boolean value and set the appropriate commandmessage
	var commandmessage string
	if value {
		commandmessage = "on"
	} else {
		commandmessage = "off"
	}

	// Send the commandmessage for the Connection method to the
	// Orchestrator ORCH server and get the acknowledgment
	acknowledge, err := client.Connection(context.Background(), &pb.Message{Message: commandmessage})

	// Check for errors and return the appropriate acknowledgement and error if any.
	if err != nil {
		return false, fmt.Errorf("call to ORCH Connection runtime failed - %v", err)
	} else {
		return acknowledge.GetSuccess(), nil
	}
}

// A function that calls the 'Observe' method of the ORCH server over a gRPC connection.
// Requires the ORCH client object and returns the stream client for the Observe service.
func Call_ORCH_Observe(client pb.OrchestratorClient) (pb.Orchestrator_ObserveClient, error) {
	// Send the valid initiation code as a Message to the Observe method of the ORCH Server
	stream, err := client.Observe(context.Background(), &pb.Message{Message: "start-stream-observe"})
	if err != nil {
		return nil, fmt.Errorf("call to ORCH Observe runtime failed - %v", err)
	}

	// Return the stream handling client for the Observe method.
	return stream, nil
}

// A function that calls the 'Status' method of the ORCH server over a gRPC connection.
// Requires the ORCH client object and returns a MeshStatus object.
func Call_ORCH_Status(client pb.OrchestratorClient) (*pb.MeshStatus, error) {
	// Call the Status method with an arbitrary message
	meshstatus, err := client.Status(context.Background(), &pb.Message{Message: "status-request"})
	if err != nil {
		return &pb.MeshStatus{}, fmt.Errorf("call to ORCH Status runtime failed - %v", err)
	}

	// Return the mesh status as a MeshStatus object.
	return meshstatus, nil
}

// A function that calls the 'Ping' method of the ORCH server over a gRPC connection.
// Requires the ORCH client object and returns a boolean success acknowledgment.
func Call_ORCH_Ping(client pb.OrchestratorClient) (bool, error) {
	// Call the Ping method with the trigger code for mesh pinging
	acknowledge, err := client.Ping(context.Background(), &pb.Message{Message: "send-ping-mesh"})

	// Check for errors and return the appropriate acknowledgement and error if any.
	if err != nil {
		return false, fmt.Errorf("call to ORCH Ping runtime failed - %v", err)
	} else {
		return acknowledge.GetSuccess(), nil
	}
}
