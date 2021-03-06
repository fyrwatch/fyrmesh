/*
===========================================================================
MIT License

Copyright (c) 2021 Manish Meganathan, Mariyam A.Ghani

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
	var triggermessage string
	if value {
		triggermessage = "setconnection-on"
	} else {
		triggermessage = "setconnection-off"
	}

	// Send the commandmessage for the Connection method to the
	// Orchestrator ORCH server and get the acknowledgment
	acknowledge, err := client.Connection(context.Background(), &pb.Trigger{Triggermessage: triggermessage})

	// Check for errors and return the appropriate acknowledgement and error if any.
	if err != nil {
		return false, fmt.Errorf("call to ORCH Connection runtime failed - %v", err)
	}

	if success := acknowledge.GetSuccess(); success {
		return true, nil
	} else {
		return false, fmt.Errorf("call to ORCH Connection returned a false acknowledge - %v", acknowledge.GetError())
	}
}

// A function that calls the 'Observe' method of the ORCH server over a gRPC connection.
// Requires the ORCH client object and returns the stream client for the Observe service.
func Call_ORCH_Observe(client pb.OrchestratorClient, sourcefilter string, typefilter string) (pb.Orchestrator_ObserveClient, error) {
	// Send the valid initiation code as a Trigger to the Observe method of the ORCH Server
	metadata := map[string]string{"source": sourcefilter, "type": typefilter}
	stream, err := client.Observe(context.Background(), &pb.Trigger{Triggermessage: "start-stream-observe", Metadata: metadata})
	if err != nil {
		return nil, fmt.Errorf("call to ORCH Observe runtime failed - %v", err)
	}

	// Return the stream handling client for the Observe method.
	return stream, nil
}

// A function that calls the 'Status' method of the ORCH server over a gRPC connection.
// Requires the ORCH client object and returns a MeshStatus object.
func Call_ORCH_Status(client pb.OrchestratorClient) (*pb.MeshOrchStatus, error) {
	// Call the Status method with an arbitrary message
	meshstatus, err := client.Status(context.Background(), &pb.Trigger{Triggermessage: "status-request"})
	if err != nil {
		return &pb.MeshOrchStatus{}, fmt.Errorf("call to ORCH Status runtime failed - %v", err)
	}

	// Return the mesh status as a MeshStatus object.
	return meshstatus, nil
}

// A function that calls the 'Ping' method of the ORCH server over a gRPC connection.
// Requires the ORCH client object and returns a boolean success acknowledgment.
func Call_ORCH_Ping(client pb.OrchestratorClient, trigger string, node string, phrase string) (bool, error) {
	// Call the Ping method with the trigger code for mesh pinging
	acknowledge, err := client.Ping(context.Background(), &pb.Trigger{Triggermessage: trigger, Metadata: map[string]string{"node": node, "phrase": phrase}})

	// Check for errors and return the appropriate acknowledgement and error if any.
	if err != nil {
		return false, fmt.Errorf("call to ORCH Ping runtime failed - %v", err)
	}

	if success := acknowledge.GetSuccess(); success {
		return true, nil
	} else {
		return false, fmt.Errorf("call to ORCH Ping returned a false acknowledge - %v", acknowledge.GetError())
	}
}

// A function that calls the 'Command' method of the ORCH server over a gRPC connection.
// Requires the ORCH client object and a command map and returns a boolean success acknowledgment.
func Call_ORCH_Command(client pb.OrchestratorClient, command map[string]string) (bool, error) {
	commandmessage := command["command"]
	delete(command, "command")

	// Call the Command method with the ControlCommand proto
	acknowledge, err := client.Command(context.Background(), &pb.ControlCommand{Command: commandmessage, Metadata: command})

	// Check for errors and return the appropriate acknowledgement and error if any.
	if err != nil {
		return false, fmt.Errorf("call to ORCH Command runtime failed - %v", err)
	}

	if success := acknowledge.GetSuccess(); success {
		return true, nil
	} else {
		return false, fmt.Errorf("call to ORCH Command returned a false acknowledge - %v", acknowledge.GetError())
	}
}

// A function that calls the 'Nodelist' method of the ORCH server over a gRPC connection.
// Requires the ORCH client and returns a slice of int64.
func Call_ORCH_Nodelist(client pb.OrchestratorClient) (map[int64]string, error) {
	// Call the Nodelist method with the Trigger proto
	nodelist, err := client.Nodelist(context.Background(), &pb.Trigger{Triggermessage: "nodelist-request"})

	// Check for errors and return the appropriate acknowledgement and error if any.
	if err != nil {
		return nil, fmt.Errorf("call to ORCH Nodelist runtime failed - %v", err)
	}

	// Obtain the slice of node IDs and return it.
	nodes := nodelist.GetNodes()
	return nodes, nil
}

// A function that calls the 'SchedulerToggle' method of the ORCH server over a gRPC connection.
// Requires the ORCH client and bool representing the toggle state.
func Call_ORCH_SchedulerToggle(client pb.OrchestratorClient, toggle bool) error {
	// Declare a trigger string
	var trigger string

	// Check the toggle
	switch toggle {
	case true:
		trigger = "setscheduler-on"
	case false:
		trigger = "setscheduler-off"
	}

	// Call the SchedulerToggle method with the Trigger proto
	acknowledge, err := client.SchedulerToggle(context.Background(), &pb.Trigger{Triggermessage: trigger})

	// Check for errors and return the appropriate acknowledgement and error if any.
	if err != nil {
		return fmt.Errorf("call to ORCH Command runtime failed - %v", err)
	}

	if success := acknowledge.GetSuccess(); success {
		return nil
	} else {
		return fmt.Errorf("call to ORCH Command returned a false acknowledge - %v", acknowledge.GetError())
	}
}

// A function that calls the 'Simulate' method of the ORCH server over a gRPC connection.
// Requires the ORCH client.
func Call_ORCH_Simulate(client pb.OrchestratorClient) error {
	// Call the SchedulerToggle method with the Trigger proto
	acknowledge, err := client.Simulate(context.Background(), &pb.Trigger{Triggermessage: "start-fire-event"})

	// Check for errors and return the appropriate acknowledgement and error if any.
	if err != nil {
		return fmt.Errorf("call to ORCH Command runtime failed - %v", err)
	}

	if success := acknowledge.GetSuccess(); success {
		return nil
	} else {
		return fmt.Errorf("call to ORCH Command returned a false acknowledge - %v", acknowledge.GetError())
	}
}
