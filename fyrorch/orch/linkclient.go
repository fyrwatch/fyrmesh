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
	"io"
	"time"

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
		logmessage = tools.NewOrchProtolog("(failure) method call failed", "LINK", "Write", err)
	} else {
		msg := fmt.Sprintf("(success) method call complete | command - %v | success - %v", commandmessage, acknowledge.GetSuccess())
		logmessage = tools.NewOrchProtolog(msg, "LINK", "Write", fmt.Errorf("%v", acknowledge.GetError()))
	}

	// Send logmessage onto the logqueue channel
	logqueue <- logmessage
}

// A function that calls the 'Read' method of the LINK server over a gRPC connection.
// Requires the LINK client object and a logqueue channel. InterfaceLogs recieved from
// LINK server will continously parsed and passed into the logqueue channel to be handled.
func Call_LINK_Read(client pb.InterfaceClient, logqueue chan tools.Log) {
	// Sleep to let other services initialize
	time.Sleep(time.Second * 5)

	// Call the 'Read' method of the LINK client with the appropriate trigger message
	stream, err := client.Read(context.Background(), &pb.Trigger{Triggermessage: "start-stream-read"})
	if err != nil {
		// Check for an error and push the protolog into the channel
		logqueue <- tools.NewOrchProtolog("(failure) method call failed.", "LINK", "Read", err)
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
			logqueue <- tools.NewOrchProtolog("(failure) method runtime failed while streaming", "LINK", "Read", errmsg)
			break
		}

		// Push the ComplexLog from the LINK into the logqueue.
		logqueue <- complexlog
	}
}
