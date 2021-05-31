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
