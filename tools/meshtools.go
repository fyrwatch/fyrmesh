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
	"strings"
)

// A function that compares if two integer slices are equal regardless of order.
// The algorithm is adopted from the StackOverflow post @ https://stackoverflow.com/a/36000696
func checkSliceEquality(x, y []int64) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of int64 -> int
	diff := make(map[int64]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}

// A struct that defines a sensor node
// and its hardware configuration values.
type SensorNode struct {
	// The identifier of the node
	NodeID int64 `firestore:"nodeID"`

	// The serial baud rate of the node
	SerialBaud int `firestore:"serialbaud"`

	// The type of DHT sensor attached
	DHTtype int `firestore:"dht_type"`

	// The pin on which the DHT sensor is attached
	DHTpin int `firestore:"dht_pin"`

	// The type of FLM sensor attached
	FLMtype int `firestore:"flm_type"`

	// The pin on which the FLM sensor is attached
	FLMpin int `firestore:"flm_pin"`

	// The type of GAS sensor attached
	GAStype int `firestore:"gas_type"`

	// The pin on which the GAS sensor is attached
	GASpin int `firestore:"gas_pin"`

	// The bool indicating if the node has a pinger button
	Pinger bool `firestore:"pinger"`

	// The pin on which the pinger button is attached
	Pingerpin int `firestore:"pinger_pin"`

	// The pin on which the connect LED is attached
	Connectpin int `firestore:"connect_pin"`
}

// A method of SensorNode that returns the
// sensor config of the node as a string.
func (sensornode *SensorNode) GetConfigString() string {
	// Declare a new slice of strings
	var configstrings []string

	if sensornode.DHTtype > 0 {
		// If DHT is set, append it to config
		configstrings = append(configstrings, "DHT")
	}

	if sensornode.FLMtype > 0 {
		// If FLM is set, append it to config
		configstrings = append(configstrings, "FLM")
	}

	if sensornode.GAStype > 0 {
		// If GAS is set, append it to config
		configstrings = append(configstrings, "GAS")
	}

	// Merge the configstrings into a single string and return it
	config := strings.Join(configstrings, "-")
	return config
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
	sensornode.NodeID, _ = strconv.ParseInt(logconfig["NODEID"], 0, 64)
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
	NodeID int64 `firestore:"nodeID"`

	// The serial baud rate of the node
	SerialBaud int `firestore:"serialbaud"`

	// The bool indicating if the node has a pinger button
	Pinger bool `firestore:"pinger"`

	// The pin on which the pinger button is attached
	Pingerpin int `firestore:"pinger_pin"`

	// The pin on which the connect LED is attached
	Connectpin int `firestore:"connect_pin"`

	// The SSID of the mesh network
	MeshSSID string `firestore:"MESHSSID"`

	// The password of the mesh network
	MeshPSWD string `firestore:"MESHPSWD"`

	// The port which the mesh network nodes communicate
	MeshPORT int `firestore:"MESHPORT"`
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
	controlnode.NodeID, _ = strconv.ParseInt(logconfig["NODEID"], 0, 64)
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

	// A map of int keys and SensorNode values that contains the list of sensor nodes on the mesh
	Nodelist map[int64]SensorNode

	// A slice of int that contains the list of all node IDs on the mesh
	NodeIDlist []int64

	// A map of string keys and MeshPing values. Maps the Pings to their respective pingIDs
	Accumulation map[string]MeshPing

	// A CloudInterface object for the mesh to connect to the database.
	Cloudinterface CloudInterface

	// A MeshDocument object that represents the mesh.
	MeshDoc MeshDocument

	// A channel of Logs that is used by all components to communicate between each other and to the console
	LogQueue chan Log

	// A channel of ObserverLogs thats is streamed to observers of the orchestrator
	ObserverQueue chan ObserverLog

	// A channel of string maps that are used to send commands to the control node
	CommandQueue chan map[string]string

	// A channel of SensorPings that are used to accumulate MeshPings
	AccumulatorQueue chan SensorPing
}

// A constructor function that generates and returns a MeshOrchestrator.
// All the channels are made with the 'make' function.
// The value of MeshConnected is false by default.
// The value of the Controlnode is set to null ControlNode until it is updated.
// The value of the ControllerID is retrieved from the config file's DeviceID.
// The value of the NodeList and NodeIDlist are set as empty slices.
func NewMeshOrchestrator() (*MeshOrchestrator, error) {
	// Create a null MeshOrchestrator
	meshorchestrator := MeshOrchestrator{}

	// Read the config file
	meshconfig, err := ReadConfig()
	if err != nil {
		return nil, fmt.Errorf("could not read config file - %v", err)
	}

	// Construct a new CloudInterface for the deviceID
	cloudinterface, err := NewCloudInterface(meshconfig.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("could not contruct cloud interface - %v", err)
	}

	// Set connection state to false by default
	meshorchestrator.MeshConnected = false
	// Set the ControllerID to the DeviceID from the config
	meshorchestrator.ControllerID = meshconfig.DeviceID
	// Set the control node of the mesh
	meshorchestrator.Controlnode = ControlNode{}
	// Set the cloud interface object to the newly generated interface
	meshorchestrator.Cloudinterface = *cloudinterface

	// Set the list of node IDs on the mesh to an emtpy slice of int
	meshorchestrator.NodeIDlist = make([]int64, 0)
	// Set the list of nodes on the mesh to an empty map of int -> SensorNode
	meshorchestrator.Nodelist = make(map[int64]SensorNode)
	// Set the list of accumulating pings to an empty map of string -> MeshPing
	meshorchestrator.Accumulation = make(map[string]MeshPing)

	// Create a log channel that will be used to pass all logs within the server.
	meshorchestrator.LogQueue = make(chan Log)
	// Create an observer channel that will be used to pass observation logs.
	meshorchestrator.ObserverQueue = make(chan ObserverLog)
	// Create a command queue that will be passed into the Orchestrator to siphon commands to the LINK.
	meshorchestrator.CommandQueue = make(chan map[string]string)
	// Create an accumulator queue that will be passed into the PingHandler to collect pings.
	meshorchestrator.AccumulatorQueue = make(chan SensorPing)

	// Create and assign the current state of the MeshOrchestrator to the MeshDoc
	meshorchestrator.MeshDoc = *NewMeshDocument(&meshorchestrator)

	// Return the meshorchestrator and nill error
	return &meshorchestrator, nil
}

// A method of MeshOrchestrator that closes all the channels and clients within it.
func (meshorchestrator *MeshOrchestrator) Close() {
	// Close the CloudInterface client
	meshorchestrator.Cloudinterface.FirestoreClient.Close()

	// Close all the channels within the MeshOrchestrator
	close(meshorchestrator.AccumulatorQueue)
	close(meshorchestrator.ObserverQueue)
	close(meshorchestrator.CommandQueue)
	close(meshorchestrator.LogQueue)
}

// A method of MeshOrchestrator that sends commands to the commandqueue
// that will trigger events that configure the Controlnode, the NodeList
// and the NodeListID fields
func (meshorchestrator *MeshOrchestrator) Initialize() {
	// Send the command to read the control node config to the CommandQueue
	command := map[string]string{"command": "readconfig-control"}
	meshorchestrator.CommandQueue <- command

	// Send the command to read the mesh node list to the CommandQueue
	command = map[string]string{"command": "readnodelist-control"}
	meshorchestrator.CommandQueue <- command
}

// A method of MeshOrchestrator that flushes the MeshDoc field to the cloud.
func (meshorchestrator *MeshOrchestrator) Flush() {
	// Update the MeshDoc with the current state of the MeshOrchestrator
	meshorchestrator.MeshDoc = *NewMeshDocument(meshorchestrator)

	// Push the meshdoc to the cloud and check the status
	err := meshorchestrator.MeshDoc.Push(&meshorchestrator.Cloudinterface)
	if err != nil {
		// Log the meshdoc failing to be flushed to the cloud.
		logmessage := NewOrchCloudlog(fmt.Sprintf("mesh document was unable to be flushed. pingID - %v", meshorchestrator.MeshDoc.ControllerID))
		meshorchestrator.LogQueue <- logmessage
	}

	// Log the meshdoc succesfully being flushed to the cloud.
	logmessage := NewOrchCloudlog(fmt.Sprintf("mesh document was flushed. pingID - %v", meshorchestrator.MeshDoc.ControllerID))
	meshorchestrator.LogQueue <- logmessage
}

// A method of MeshOrchestrator that accepts a Log of type 'nodelist' and parses the nodelist sequence
// on it into a slice of integer nodeIDs and assigns it to the NodeIDlist field. Finally calls the
// method to update the NodeList based on the new NodeIDlist.
func (meshorchestrator *MeshOrchestrator) SetNodeIDlist(log Log) error {
	// Retrieve the type of the Log
	logtype := log.GetLogtype()
	// Check if the logtype is 'controlconfig'
	if logtype != "nodelist" {
		return fmt.Errorf("log is not of type 'nodelist'")
	}

	// Retrieve the log metadata
	logmetadata := log.GetLogmetadata()
	// Retrieve the nodelist sequence from the metadata
	seqnodelist := logmetadata["nodelist"]
	// Trim the nodelist sequence for trailing splitters
	seqnodelist = strings.TrimSuffix(seqnodelist, "-")
	// Split the nodelist sequence into a slice of strings
	strnodelist := strings.Split(seqnodelist, "-")

	// Declare nodelist of type slice of int
	var nodelist []int64
	// Iterate over the string nodelist slice
	for _, strnode := range strnodelist {
		// Convert the string to an int and append it to the int nodelist
		node, _ := strconv.ParseInt(strnode, 0, 64)
		nodelist = append(nodelist, node)
	}

	// Assign the new NodeIDlist
	meshorchestrator.NodeIDlist = nodelist
	// Call the method to update the NodeList based on the new NodeIDlist
	meshorchestrator.UpdateNodelist()
	// Call the method to update the MeshDocument and flush it
	meshorchestrator.Flush()
	return nil
}

// A method of MeshOrchestrator that sets the value of the Controlnode field.
// Accepts a log of type 'controlconfig' and constructs a ControlNode from the log.
func (meshorchestrator *MeshOrchestrator) SetControlnode(log Log) error {
	// Retrieve the type of the Log
	logtype := log.GetLogtype()
	// Check if the logtype is 'controlconfig'
	if logtype != "controlconfig" {
		return fmt.Errorf("log is not of type 'controlconfig'")
	}

	// Construct a new ControlNode
	controlnode, err := NewControlNode(log)
	if err != nil {
		return fmt.Errorf("control node config could not be constructed - %v", err)
	}

	// Assign the controlnode to the meshorchestrator
	meshorchestrator.Controlnode = *controlnode
	// Call the method to update the MeshDocument and flush it
	meshorchestrator.Flush()
	return nil
}

// A method of MeshOrchestrator that adds/updates the Node on the Nodelist map.
// Accepts a log og type 'configdata' and constructs a SensorNode from the log.
// The SensorNode is then added to the Nodelist with the NodeID being the key.
func (meshorchestrator *MeshOrchestrator) SetNode(log Log) error {
	// Retrieve the type of the Log
	logtype := log.GetLogtype()
	// Check if the logtype is 'controlconfig'
	if logtype != "configdata" {
		return fmt.Errorf("log is not of type 'configdata'")
	}

	// Construct a new SensorNode
	sensornode, err := NewSensorNode(log)
	if err != nil {
		return fmt.Errorf("sensor node config could not be constructed - %v", err)
	}

	// Assign the sensornode to the meshorchestrator's Nodelist
	meshorchestrator.Nodelist[sensornode.NodeID] = *sensornode
	// Call the method to update the MeshDocument and flush it
	meshorchestrator.Flush()
	return nil
}

// A method of MeshOrchestator that constructs a SensorPing object from the
// log and sends it to the AccumulatorQueue if it is not a userping
func (meshorchestrator *MeshOrchestrator) SetSensorData(log Log) error {
	// Retrieve the type of the Log
	logtype := log.GetLogtype()
	// Check if the logtype is 'controlconfig'
	if logtype != "sensordata" {
		return fmt.Errorf("log is not of type 'sensordata'")
	}

	// Construct a new SensorNode
	sensorping, err := NewSensorPing(log, meshorchestrator)
	if err != nil {
		return fmt.Errorf("sensor ping could not be constructed - %v", err)
	}

	// userpings are never accumulated
	userping := strings.HasPrefix(sensorping.PingID, "userping")
	// only mesh wide pings can be accumulated
	meshping := strings.HasSuffix(sensorping.PingID, "mesh")

	// Check if sensor ping is not a user ping and is a mesh ping
	if !userping && meshping {
		meshorchestrator.AccumulatorQueue <- *sensorping
	}

	return nil
}

// A method of MeshOrchestrator that sends the command to
// retreieve a new copy of the nodelist from the mesh.
func (meshorchestrator *MeshOrchestrator) UpdateNodeIDlist() {
	// Send the command to read a copy of the current mesh node list to the CommandQueue
	command := map[string]string{"command": "readnodelist-control"}
	meshorchestrator.CommandQueue <- command
}

// A method of MeshhOrchestrator that sets updates the NodeList field.
// Compares the NodeIDlist field with a slice of NodeID integers collected
// from the NodeList. If they are equal, no updation is performed, otherwise
// a command is sent to ping the entire mesh for configdata, each of which
// will accumulate into the NodeList map.
func (meshorchestrator *MeshOrchestrator) UpdateNodelist() {
	// Retrieve the current Nodelist
	oldNodelist := meshorchestrator.Nodelist
	// Retrieve the current(updated) NodeIDlist
	newNodeIDlist := meshorchestrator.NodeIDlist

	// Declare a slice of int
	var oldNodeIDlist []int64
	// Iterate over the keys of the Nodelits
	for nodeid := range oldNodelist {
		// Append the integer keys into the slice
		oldNodeIDlist = append(oldNodeIDlist, nodeid)
	}

	// Check if the two NodeIDlists are equal
	result := checkSliceEquality(oldNodeIDlist, newNodeIDlist)
	if !result {
		// If they are not equal, send the command to ping the mesh for config data to the CommandQueue
		command := map[string]string{"command": "readconfig-mesh", "ping": fmt.Sprintf("controlping-nodelistupdater-%v-mesh", CurrentISOtime())}
		meshorchestrator.CommandQueue <- command
	}
}

// A method of MeshOrchestrator that returns a simplified nodelist.
// The simplified nodelist is mapping of the nodeIDs to their config strings.
func (meshorchestrator *MeshOrchestrator) GetSimpleNodeList() map[int64]string {
	// Create a null map and retrieve the NodeList from the meshorchestrator
	simplenodelist := make(map[int64]string)
	nodelist := meshorchestrator.Nodelist

	// Iterate over the nodelist and accumulate into the simple map
	for nodeid, node := range nodelist {
		simplenodelist[nodeid] = node.GetConfigString()
	}

	// Return the simple nodelist
	return simplenodelist
}
