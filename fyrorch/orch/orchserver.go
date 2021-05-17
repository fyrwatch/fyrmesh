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
	"net"

	pb "github.com/fyrwatch/fyrmesh/proto"
	tools "github.com/fyrwatch/fyrmesh/tools"
	"google.golang.org/grpc"
)

// A struct that defines the Orchestrator gRPC Server
type OrchestratorServer struct {
	pb.UnimplementedOrchestratorServer
	meshconnected  bool
	meshidentifier string
	commandqueue   chan string
	observerqueue  chan string
	logqueue       chan string
}

// A function that implements the 'Connection' method of the Orchestrator service.
// Accepts a Message and returns an Acknowledge
func (server *OrchestratorServer) Connection(ctx context.Context, trigger *pb.Trigger) (*pb.Acknowledge, error) {
	// Retrieve the trigger message from the Message object
	triggermessage := trigger.GetTriggermessage()

	// Check the value of the command message
	switch triggermessage {
	case "setconnection-on":
		// Set the meshconnected value to True
		server.meshconnected = true
		// Send a command to the server's command queue
		server.commandqueue <- "connection-on"

	case "setconnection-off":
		// Set the meshconnected value to True
		server.meshconnected = false
		// Send a command to the server's command queue
		server.commandqueue <- "connection-off"

	default:
		// Default to returning a fail Acknowledge because of an unsupported command message
		return &pb.Acknowledge{Success: false, Error: "unsupported command"}, nil
	}

	// Return an success Acknowledge with no error
	return &pb.Acknowledge{Success: true, Error: "nil"}, nil
}

// A function that implements the 'Observe' method of the Orchestrator service.
// Accepts a Message and returns a stream of Message
func (server *OrchestratorServer) Observe(trigger *pb.Trigger, stream pb.Orchestrator_ObserveServer) error {
	// Retrieve the stream trigger message from the Message object and check its value.
	triggermessage := trigger.GetTriggermessage()
	if triggermessage != "start-stream-observe" {
		// If stream initiation code is invalid. Send one error message over the stream and return.
		stream.Send(&pb.SimpleLog{Message: "invalid observe stream initiation code"})
		return nil
	}

	// Send the signal to enable the observer queue for the log handler.
	server.logqueue <- "enable-observe"

	// Iterate over the observer channel
	for log := range server.observerqueue {
		// Send each log recieved on the channel to the stream.
		err := stream.Send(&pb.SimpleLog{Message: log})
		if err != nil {
			return err
		}
	}

	// Send the signal to disable the observer queue for the log handler.
	server.logqueue <- "disable-observe"
	return nil
}

// A function that implements the 'Status' method of the Orchestrator service.
// Accepts a Message and returns a MeshStatus
func (server *OrchestratorServer) Status(ctx context.Context, trigger *pb.Trigger) (*pb.MeshStatus, error) {
	// Return values from the server configuration as MeshStatus object.
	return &pb.MeshStatus{
		MeshID:    server.meshidentifier,
		Connected: server.meshconnected,
	}, nil
}

// A function that implements the 'Ping' method of the Orchestrator service.
// Accepts a Message and returns an Acknowledge
func (server *OrchestratorServer) Ping(ctx context.Context, trigger *pb.Trigger) (*pb.Acknowledge, error) {
	// Retrieve the trigger message from the Message object
	triggermessage := trigger.GetTriggermessage()

	// Check the value of the command message
	switch triggermessage {
	case "send-ping-mesh":
		// Send a command to the server's command queue
		server.commandqueue <- "readsensors-mesh"

	case "send-ping-node":
		// Returning a fail Acknowledge because of an unimplemented command message.
		return &pb.Acknowledge{Success: false, Error: "unimplemented command"}, nil

	default:
		// Default to returning a fail Acknowledge because of an unsupported command message.
		return &pb.Acknowledge{Success: false, Error: "unsupported command"}, nil
	}

	// Return an success Acknowledge with no error
	return &pb.Acknowledge{Success: true, Error: "nil"}, nil
}

// A function that handles the output of the commands recieved over a given command queue
// by passing each recieved command to function that calls the the 'Write' method of the
// interface LINK server. Iterates infinitely until the commandqueue is closed.
func pushcommands(linkclient pb.InterfaceClient, logqueue chan string, commandqueue chan string) {
	for command := range commandqueue {
		Call_LINK_Write(linkclient, logqueue, command)
	}
}

// A function that creates the gRPC server for the Orchestrator ORCH service
// and sets it to listen on the appropriate port. Starts a go routine to check
// the server's command queue
func Start_ORCH_Server(linkclient pb.InterfaceClient, logqueue chan string, commandqueue chan string, obsqueue chan string) error {
	// Read the config file
	config, err := tools.ReadConfig()
	if err != nil {
		return fmt.Errorf("could not read config file - %v", err)
	}

	// Construct the port string from the Port field in the ORCH ServiceConfig within the Config
	port := fmt.Sprintf(":%d", config.Services["ORCH"].Port)
	// Setup the listener on the constructed port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("could not set up listener on the port tcp%v - %v", port, err)
	}

	// Create the gRPC server and register it with the commandqueue, observerqueue and the logqueue
	grpcserver := grpc.NewServer()
	pb.RegisterOrchestratorServer(grpcserver, &OrchestratorServer{
		meshconnected:  false,
		meshidentifier: config.DeviceID,
		commandqueue:   commandqueue,
		observerqueue:  obsqueue,
		logqueue:       logqueue,
	})

	// Start a go-routine to check the server's command queue and push them to LINK server
	go pushcommands(linkclient, logqueue, commandqueue)

	// Serve the gRPC server on the listener port
	if err := grpcserver.Serve(listener); err != nil {
		return fmt.Errorf("could not start the server - %v", err)
	} else {
		return nil
	}
}
