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
	"time"

	"google.golang.org/grpc"

	pb "github.com/fyrwatch/fyrmesh/proto"
	tools "github.com/fyrwatch/fyrmesh/tools"
)

// A struct that defines the Orchestrator gRPC Server
type OrchestratorServer struct {
	pb.UnimplementedOrchestratorServer
	meshorchestrator *tools.MeshOrchestrator
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
		server.meshorchestrator.MeshConnected = true
		// Send a command to the server's command queue
		command := map[string]string{"command": "connection-on"}
		server.meshorchestrator.CommandQueue <- command

	case "setconnection-off":
		// Set the meshconnected value to True
		server.meshorchestrator.MeshConnected = false
		// Send a command to the server's command queue
		command := map[string]string{"command": "connection-off"}
		server.meshorchestrator.CommandQueue <- command

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
	triggermetadata := trigger.GetMetadata()
	// Retrieve the observe filters from the Trigger metadata field.
	typefilter := triggermetadata["type"]
	sourcefilter := triggermetadata["source"]

	if triggermessage != "start-stream-observe" {
		// If stream initiation code is invalid. Send one error message over the stream and return.
		stream.Send(&pb.SimpleLog{Message: "[OBS] invalid observe stream initiation code"})
		return nil
	}

	// Send the signal to enable the observer queue for the log handler.
	obstoggle := tools.NewObsCommand("enable-observe")
	server.meshorchestrator.LogQueue <- obstoggle

	// Iterate over the observer channel
	for log := range server.meshorchestrator.ObserverQueue {

		if sourcefilter == "" && typefilter == "" {
			// If both filters are not set - send all logs recieved on the channel to the stream.
			if err := stream.Send(&pb.SimpleLog{Message: log.Logmessage}); err != nil {
				return err
			}

		} else if sourcefilter != "" && typefilter == "" {
			// If only source filter is set - check the source of logs recieved on the channel and send on the stream
			if log.Logsource == sourcefilter {
				if err := stream.Send(&pb.SimpleLog{Message: log.Logmessage}); err != nil {
					return err
				}
			}

		} else if sourcefilter == "" && typefilter != "" {
			// If only type filter is set - check the type of logs recieved on the channel and send on the stream
			if log.Logsource == typefilter {
				if err := stream.Send(&pb.SimpleLog{Message: log.Logmessage}); err != nil {
					return err
				}
			}

		} else if sourcefilter != "" && typefilter != "" {
			// If both filters are set - check the type and source of logs receieved on the channel and send on the stream
			if log.Logsource == sourcefilter && log.Logtype == typefilter {
				if err := stream.Send(&pb.SimpleLog{Message: log.Logmessage}); err != nil {
					return err
				}
			}
		}
	}

	// Send the signal to disable the observer queue for the log handler.
	obstoggle = tools.NewObsCommand("disable-observe")
	server.meshorchestrator.LogQueue <- obstoggle

	return nil
}

// A function that implements the 'Status' method of the Orchestrator service.
// Accepts a Message and returns a MeshStatus
func (server *OrchestratorServer) Status(ctx context.Context, trigger *pb.Trigger) (*pb.MeshOrchStatus, error) {
	// Return values from the server configuration as a MeshOrchStatus object.
	return &pb.MeshOrchStatus{
		Connected:     server.meshorchestrator.MeshConnected,
		ControllerID:  server.meshorchestrator.ControllerID,
		ControlnodeID: int64(server.meshorchestrator.Controlnode.NodeID),
		Nodelist:      &pb.NodeList{Nodes: server.meshorchestrator.GetSimpleNodeList()},
		MeshSSID:      server.meshorchestrator.Controlnode.MeshSSID,
		MeshPSWD:      server.meshorchestrator.Controlnode.MeshPSWD,
		MeshPORT:      int32(server.meshorchestrator.Controlnode.MeshPORT),
	}, nil
}

// A function that implements the 'Ping' method of the Orchestrator service.
// Accepts a Message and returns an Acknowledge
func (server *OrchestratorServer) Ping(ctx context.Context, trigger *pb.Trigger) (*pb.Acknowledge, error) {
	// Retrieve the trigger message and metadata from the Trigger proto
	triggermessage := trigger.GetTriggermessage()
	triggermetadata := trigger.GetMetadata()

	// Check the value of the trigger message
	switch triggermessage {
	case "ping-sensor-mesh":
		// Send a command to ping mesh for sensor data to the server's command queue
		command := map[string]string{"command": "readsensors-mesh", "ping": fmt.Sprintf("userping-%v-mesh", triggermetadata["phrase"])}
		server.meshorchestrator.CommandQueue <- command

	case "ping-sensor-node":
		// Send a command to ping a node for sensor data to the server's command queue
		command := map[string]string{"command": "readsensors-node", "ping": fmt.Sprintf("userping-%v-node", triggermetadata["phrase"]), "node": triggermetadata["node"]}
		server.meshorchestrator.CommandQueue <- command

	case "ping-config-mesh":
		// Send a command to ping mesh for config data to the server's command queue
		command := map[string]string{"command": "readconfig-mesh", "ping": fmt.Sprintf("userping-%v-mesh", triggermetadata["phrase"])}
		server.meshorchestrator.CommandQueue <- command

	case "ping-config-node":
		// Send a command to ping a node for config data to the server's command queue
		command := map[string]string{"command": "readconfig-node", "ping": fmt.Sprintf("userping-%v-node", triggermetadata["phrase"]), "node": triggermetadata["node"]}
		server.meshorchestrator.CommandQueue <- command

	case "ping-control":
		// Send a command to ping the control node for config data to the server's command queue
		command := map[string]string{"command": "readconfig-control"}
		server.meshorchestrator.CommandQueue <- command

	default:
		// Default to returning a fail Acknowledge because of an unsupported command message.
		return &pb.Acknowledge{Success: false, Error: "unsupported trigger command"}, nil
	}

	// Return an success Acknowledge with no error
	return &pb.Acknowledge{Success: true, Error: "nil"}, nil
}

// A function that implements the 'Command' method of the Orchestrator service.
// Accepts a ControlCommand and returns an Acknowledge
func (server *OrchestratorServer) Command(ctx context.Context, controlcommand *pb.ControlCommand) (*pb.Acknowledge, error) {
	// Retrieve the command message and metadata from the ControlCommand object
	commandmessage := controlcommand.GetCommand()
	commandmetadata := controlcommand.GetMetadata()

	// Set the command message as the 'command' key in a new map
	command := map[string]string{"command": commandmessage}
	// Collect the metadata values into the same command map
	for key, value := range commandmetadata {
		command[key] = value
	}

	// Send the command over the CommandQueue
	server.meshorchestrator.CommandQueue <- command
	// Return an success Acknowledge with no error
	return &pb.Acknowledge{Success: true, Error: "nil"}, nil
}

// A function that implements the 'Nodelist' method of the Orchestrator service.
// Accepts a Trigger and returns a NodeList.
func (server *OrchestratorServer) Nodelist(ctx context.Context, trigger *pb.Trigger) (*pb.NodeList, error) {
	// Return the list of nodes currently on the mesh as a NodeList proto
	return &pb.NodeList{Nodes: server.meshorchestrator.GetSimpleNodeList()}, nil
}

// A function that handles the output of the commands recieved over a given command queue
// by passing each recieved command to function that calls the the 'Write' method of the
// interface LINK server. Iterates infinitely until the commandqueue is closed.
func CommandHandler(linkclient pb.InterfaceClient, meshorchestrator *tools.MeshOrchestrator) {
	for command := range meshorchestrator.CommandQueue {
		Call_LINK_Write(linkclient, meshorchestrator.LogQueue, command)
	}
}

// A function that handles the scheduled pinging of the message at a regualar interval
// The interval is defined in the config file as an integer number of seconds.
// The scheduler waits 15s before starting the pings to give time for the mesh and
// orchestrator to initialize when the service first starts.
func Scheduler(meshorchestrator *tools.MeshOrchestrator, pingrate int) {
	// Sleep for 15s to give time for other orchestrator services to initialize
	time.Sleep(time.Second * 15)
	// Log the beginning of the scheduled pinging to the LogQueue
	logmessage := tools.NewOrchSchedlog(fmt.Sprintf("ping scheduler has started with ping rate %v", pingrate))
	meshorchestrator.LogQueue <- logmessage

	for {
		// Generate a ping ID and command to ping the mesh for sensors and push it to the commandQueue
		pingid := fmt.Sprintf("controlping-scheduler-%v-mesh", tools.CurrentISOtime())
		command := map[string]string{"command": "readsensors-mesh", "ping": pingid}
		meshorchestrator.CommandQueue <- command

		// Log the scheduled ping with the ping ID.
		logmessage := tools.NewOrchSchedlog(fmt.Sprintf("mesh pinged for sensor data. pingID -  %v", pingid))
		meshorchestrator.LogQueue <- logmessage
		// Sleep for the pingrate number of seconds.
		time.Sleep(time.Second * time.Duration(pingrate))
	}
}

// A function that creates the gRPC server for the Orchestrator ORCH service
// and sets it to listen on the appropriate port. Starts a go routine to check
// the server's command queue
func Start_ORCH_Server(linkclient pb.InterfaceClient, meshorchestrator *tools.MeshOrchestrator) error {
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
	pb.RegisterOrchestratorServer(grpcserver, &OrchestratorServer{meshorchestrator: meshorchestrator})

	// Start a go-routine to check the server's command queue and push them to LINK server.
	go CommandHandler(linkclient, meshorchestrator)

	// Start a go-routine to check the servers's accumulation queue and handle the recieved pings.
	go tools.PingHandler(meshorchestrator)

	// Start a go-routine to send scheduled pings to the mesh
	go Scheduler(meshorchestrator, config.SchedulerPingRate)

	// Call the Initialize method the meshorchestrator to configure the node list and control node fields.
	meshorchestrator.Initialize()

	// Serve the gRPC server on the listener port
	if err := grpcserver.Serve(listener); err != nil {
		return fmt.Errorf("could not start the server - %v", err)
	} else {
		return nil
	}
}
