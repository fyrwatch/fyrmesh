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
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "github.com/fyrwatch/fyrmesh/proto"
	tools "github.com/fyrwatch/fyrmesh/tools"
)

// A function that establishes a gRPC connection to the interface LINK server and returns the LINK
// client object, the gPRC connection object and any error that occurs while attempting to connect.
func GRPCconnect_LINK() (*pb.InterfaceClient, *grpc.ClientConn, error) {
	// Read the service config for the LINK server
	config, err := tools.ReadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("could not read config file - %v", err)
	}

	// Construct the URL for the interface LINK server
	linkconfig := config.Services["LINK"]
	linkhost := fmt.Sprintf("%s:%d", linkconfig.Host, linkconfig.Port)

	// Connect to the Interface LINK gRPC Server
	conn, err := grpc.Dial(linkhost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, nil, fmt.Errorf("could not dialup LINK gRPC server - %v", err)
	}

	// Create an Interface LINK Client and return it along with the connection object and a nil error
	client := pb.NewInterfaceClient(conn)
	return &client, conn, nil
}

// A function that calls the 'Write' method of the LINK server over a gRPC connection.
// Requires the LINK client object, a logqueue channel and the string command to send.
func Call_LINK_Write(client pb.InterfaceClient, logqueue chan tools.Log, command map[string]string) {
	commandmessage := command["command"]
	delete(command, "command")

	// Send a string command to the Interface LINK server and get the acknowledgment
	acknowledge, err := client.Write(context.Background(), &pb.ControlCommand{Command: commandmessage, Metadata: command})

	// Check for errors and construct appropriate protolog
	var logmessage *tools.OrchLog
	if err != nil {
		logmessage = tools.NewOrchProtolog("method call failed.", "LINK", "Write", err)
	} else {
		msg := fmt.Sprintf("method call complete. command - %v. success - %v", commandmessage, acknowledge.GetSuccess())
		logmessage = tools.NewOrchProtolog(msg, "LINK", "Write", fmt.Errorf("%v", acknowledge.GetError()))
	}

	// Send logmessage onto the logqueue channel
	logqueue <- logmessage
}

// A function that calls the 'Read' method of the LINK server over a gRPC connection.
// Requires the LINK client object and a logqueue channel. InterfaceLogs recieved from
// LINK server will continously parsed and passed into the logqueue channel to be handled.
func Call_LINK_Read(client pb.InterfaceClient, logqueue chan tools.Log) {
	// Call the 'Read' method of the LINK client with the appropriate trigger message
	stream, err := client.Read(context.Background(), &pb.Trigger{Triggermessage: "start-stream-read"})
	if err != nil {
		// Check for an error and push the protolog into the channel
		logmessage := tools.NewOrchProtolog("method call failed.", "LINK", "Read", err)
		logqueue <- logmessage
	}

	// Start an infinite loop to read from the stream
	for {
		// Recieve an InterfaceLog object from the stream
		complexlog, err := stream.Recv()

		// Break out of loop if stream has closed
		if err == io.EOF {
			break
		}

		// Push to logqueue if any other error occurs and break from loop.
		if err != nil {
			errstatus, _ := status.FromError(err)
			errmsg := fmt.Errorf("StreamError - (%v)%v", errstatus.Code(), errstatus.Message())
			logmessage := tools.NewOrchProtolog("method runtime failed while streaming", "LINK", "Read", errmsg)
			logqueue <- logmessage
			break
		}

		// Push the ComplexLog from the LINK into the logqueue.
		logqueue <- complexlog
	}
}
