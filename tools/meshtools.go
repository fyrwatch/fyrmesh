/*
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.

This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed
without the express permission of Manish Meganathan and Mariyam A.Ghani
===========================================================================
FyrMesh gopkg tools
===========================================================================
*/
package tools

import (
	"fmt"
	"strconv"
)

// A struct that defines a sensor node
// and its hardware configuration values.
type SensorNode struct {
	// The identifier of the node
	NodeID int

	// The serial baud rate of the node
	SerialBaud int

	// The type of DHT sensor attached
	DHTtype int

	// The pin on which the DHT sensor is attached
	DHTpin int

	// The type of FLM sensor attached
	FLMtype int

	// The pin on which the FLM sensor is attached
	FLMpin int

	// The type of GAS sensor attached
	GAStype int

	// The pin on which the GAS sensor is attached
	GASpin int

	// The bool indicating if the node has a pinger button
	Pinger bool

	// The pin on which the pinger button is attached
	Pingerpin int

	// The pin on which the connect LED is attached
	Connectpin int
}

// A constructor function that generates and returns a SensorNode.
// Only accepts a Log of type 'configdata'.
func NewSensorNode(log Log) (*SensorNode, error) {
	// Retrieve the type of the Log
	logtype := log.GetLogtype()
	// Check if the logtype is 'configdata'
	if logtype != "configdata" {
		return nil, fmt.Errorf("log is not of type 'configdata'")
	}

	// Retrieve the metadata of the log and deep deserialize the config
	metadata := log.GetLogmetadata()
	logconfig := Deepdeserialize(metadata["config"])

	// Create a null SensorNode
	sensornode := SensorNode{}
	// Parse and assign the general hardware config values
	sensornode.NodeID, _ = strconv.Atoi(logconfig["NODEID"])
	sensornode.SerialBaud, _ = strconv.Atoi(logconfig["SERIALBAUD"])
	sensornode.Pinger, _ = strconv.ParseBool(logconfig["PINGER"])
	sensornode.Pingerpin, _ = strconv.Atoi(logconfig["PINGERPIN"])
	sensornode.Connectpin, _ = strconv.Atoi(logconfig["CONNECTLEDPIN"])
	// Parse and assign the sensor hardware config values
	sensornode.DHTtype, _ = strconv.Atoi(logconfig["DHTTYP"])
	sensornode.DHTpin, _ = strconv.Atoi(logconfig["DHTPIN"])
	sensornode.FLMtype, _ = strconv.Atoi(logconfig["FLMTYP"])
	sensornode.FLMpin, _ = strconv.Atoi(logconfig["FLMPIN"])
	sensornode.GAStype, _ = strconv.Atoi(logconfig["GASTYP"])
	sensornode.GASpin, _ = strconv.Atoi(logconfig["GASPIN"])

	// Return the pointer of the sensor node and a nil error
	return &sensornode, nil
}

// A struct that defines a control node with its hardware configuration
// values along with the configuration values that define the mesh.
type ControlNode struct {
	// The identifier of the node
	NodeID int

	// The serial baud rate of the node
	SerialBaud int

	// The bool indicating if the node has a pinger button
	Pinger bool

	// The pin on which the pinger button is attached
	Pingerpin int

	// The pin on which the connect LED is attached
	Connectpin int

	// The SSID of the mesh network
	MeshSSID string

	// The password of the mesh network
	MeshPSWD string

	// The port which the mesh network nodes communicate
	MeshPORT int
}

// A constructor function that generates and returns a ControlNode.
// Only accepts a Log of type 'controlconfig'.
func NewControlNode(log Log) (*ControlNode, error) {
	// Retrieve the type of the Log
	logtype := log.GetLogtype()
	// Check if the logtype is 'controlconfig'
	if logtype != "controlconfig" {
		return nil, fmt.Errorf("log is not of type 'controlconfig'")
	}

	// Retrieve the metadata of the log and deep deserialize the config
	metadata := log.GetLogmetadata()
	logconfig := Deepdeserialize(metadata["config"])

	// Create a null SensorNode
	controlnode := ControlNode{}
	// Parse and assign the general hardware config values
	controlnode.NodeID, _ = strconv.Atoi(logconfig["NODEID"])
	controlnode.SerialBaud, _ = strconv.Atoi(logconfig["SERIALBAUD"])
	controlnode.Pinger, _ = strconv.ParseBool(logconfig["PINGER"])
	controlnode.Pingerpin, _ = strconv.Atoi(logconfig["PINGERPIN"])
	controlnode.Connectpin, _ = strconv.Atoi(logconfig["CONNECTLEDPIN"])
	// Parse and assign the mesh config values
	controlnode.MeshSSID = logconfig["MESH_SSID"]
	controlnode.MeshPSWD = logconfig["MESH_PSWD"]
	controlnode.MeshPORT, _ = strconv.Atoi(logconfig["MESH_PORT"])

	// Return the pointer of the control node and a nil error
	return &controlnode, nil
}

// A struct that defines a mesh orchestrator. It contains all the
// core structures that are shared among its various sub-routines.
type MeshOrchestrator struct {
	// A bool indicating if the connection state to the control node has been set
	MeshConnected bool

	// A string identifier of the controller that is running the orchestrator
	ControllerID string

	// A ControlNode object that contains the configuration of the mesh control node
	Controlnode ControlNode

	// A slice of SensorNode objects that contains the list of sensor nodes on the mesh
	Nodelist []SensorNode

	// A channel of Logs that is used by all components to communicate between each other and to the console
	LogQueue chan Log

	// A channel of ObserverLogs thats is streamed to observers of the orchestrator
	ObserverQueue chan ObserverLog

	// A channel of string maps that are used to send commands to the control node
	CommandQueue chan map[string]string
}

// A constructor function that generates and returns a MeshOrchestrator.
// All the channels are made with the 'make' function.
// The value of MeshConnected is false by default.
// The value of the Controlnode is set to null ControlNode until it is updated.
// The value of the ControllerID is retrieved from the config file's DeviceID.
func NewMeshOrchestrator() (*MeshOrchestrator, error) {
	// Create a null MeshOrchestrator
	meshorchestrator := MeshOrchestrator{}

	// Set connection state to false by default
	meshorchestrator.MeshConnected = false
	// Set the control node of the mesh
	meshorchestrator.Controlnode = ControlNode{}

	// Create a log channel that will be used to pass all logs within the server.
	meshorchestrator.LogQueue = make(chan Log)
	// Create an observer channel that will be used to pass observation logs.
	meshorchestrator.ObserverQueue = make(chan ObserverLog)
	// Create a command queue that will be passed into the Orchestrator to siphon commands to the LINK.
	meshorchestrator.CommandQueue = make(chan map[string]string)

	// Read the config file
	meshconfig, err := ReadConfig()
	if err != nil {
		return nil, fmt.Errorf("could not read config file - %v", err)
	}

	// Set the ControllerID to the DeviceID from the config
	meshorchestrator.ControllerID = meshconfig.DeviceID
	// Return the meshorchestrator and nill error
	return &meshorchestrator, nil
}

// A method of MeshOrchestrator that closes all the channels inside it.
func (meshorchestrator *MeshOrchestrator) Close() {
	// Close all the channels within the MeshOrchestrator
	close(meshorchestrator.ObserverQueue)
	close(meshorchestrator.CommandQueue)
	close(meshorchestrator.LogQueue)
}
