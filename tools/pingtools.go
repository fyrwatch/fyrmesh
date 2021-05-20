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

// A struct that defines the ping response
// of sensordata from a sensor node.
type SensorPing struct {
	// A mapping of string sensor types to float64 sensor values
	Sensordata map[string]float64

	// A SensorNode object that represents the origin of the ping data
	Sensornode SensorNode

	// A string that represents the pingID of a ping response
	PingID string

	// A string that represents the time of the ping response
	Pingtime string

	// A float32 value that reprsents the probability of fire in the neighbourhood of the node
	Fireprobability float64
}

// A function that calculates the probability of a fire in the
// neighbourhood of the node and sets it to the SensorPing object
func (sensorping *SensorPing) CalculateFireProbability() error {
	// TODO: implement
	return nil
}

// A constructor function that generates and returns a SensorPing.
// Requires a Log of type 'sensordata'
// The value of Sensordata is collected by parsed the sensordata log.
// The value of the Sensornode is retrieved from MeshOrchestrator's Nodelist.
// The value of the PingID is taken from sensordata log.
// The value of the Pingtime is the current time when the struct is constructed.
func NewSensorPing(log Log, meshorchestrator *MeshOrchestrator) (*SensorPing, error) {
	// Retrieve the type of the Log
	logtype := log.GetLogtype()
	// Check if the logtype is 'sensordata'
	if logtype != "sensordata" {
		return nil, fmt.Errorf("log is not of type 'sensordata'")
	}

	// Retrieve the metadata of the log and deep deserialize the sensors
	metadata := log.GetLogmetadata()
	sensordata := Deepdeserialize(metadata["sensors"])

	// Create an empty SensorPing
	sensorping := SensorPing{}
	sensorping.Sensordata = make(map[string]float64)

	// Parse and assign the sensor values that exist
	sensorkeys := []string{"HUM", "TEM", "FLM", "GAS"}
	for _, key := range sensorkeys {
		if val, ok := sensordata[key]; ok {
			parsedval, _ := strconv.ParseFloat(val, 64)
			sensorping.Sensordata[key] = parsedval
		}
	}

	// Retrieve the node ID from the metadata
	nodeid, _ := strconv.ParseInt(metadata["node"], 0, 64)
	// Retrieve the SensorNode object for the nodeID from the mesh orchestrator
	sensorping.Sensornode = meshorchestrator.Nodelist[nodeid]
	// Assign the ping ID from the metadata
	sensorping.PingID = metadata["ping"]
	// Assign the ping time to the current time
	sensorping.Pingtime = CurrentISOtime()

	// Calculate the value of the fire probability
	// sensorping.CalculateFireProbability()
	// Return the sensor ping
	return &sensorping, nil
}

// A struct represents a collection
// of SensorPing of the same ping ID
type MeshPing struct {
	// A slice of int64 node IDs that are expected to respond to the ping
	Nodelist []int64

	// A mapping of int64 nodeIDs to their respective SensorPings
	Pings map[int64]SensorPing

	// A string that represents the shared ping ID of all SensorPings in Pings
	PingID string

	// A string that represents the time of response of the first SensorPing to get accumulated
	Pingtime string
}

// A constructor function that generates and returns a MeshPing.
// Requires a pingID, a ping time and a slice of int64 representing the list of nodes from which to expect response.
// The value of PingID is set based on the value passed to constructor.
// The value of Ping time is set based on the value passed to constructor.
// The value of Nodelist is set based on the value passed to constructor.
// The value of the Pings is an empty map of int64 -> SensorPing
func NewMeshPing(pingid string, pingtime string, nodelist []int64) *MeshPing {
	// Create an empty MeshPing
	meshping := MeshPing{}

	// Assign the ping id, ping time and nodelist
	meshping.PingID = pingid
	meshping.Pingtime = pingtime
	meshping.Nodelist = nodelist
	// Create and assign empty slices for Pings
	meshping.Pings = make(map[int64]SensorPing)

	// Return the meshping
	return &meshping
}

// A method of MeshPing that returns a bool indicating whether
// all the nodes in its Nodelist have responded to the ping.
func (meshping *MeshPing) Complete() bool {
	// Iterate over the nodelist
	for _, nodeid := range meshping.Nodelist {
		// Check if the node ID has been added to the Pings
		if _, exists := meshping.Pings[nodeid]; !exists {
			return false
		}
	}
	return true
}

// A method of MeshPing that generates and returns a mappings of the string node ID to its Sensordata map
func (meshping *MeshPing) GenerateSensordatamap() map[string]map[string]float64 {
	// Create an empty sensordata map
	sensordata := make(map[string]map[string]float64)

	// Iterate over the Pings in the meshping
	for nodeid, sensorping := range meshping.Pings {
		// Convert the nodeIDs to strings and assign the Sensordata
		sensordata[strconv.FormatInt(nodeid, 10)] = sensorping.Sensordata
	}

	// Return the sensordata
	return sensordata
}

// A method of MeshPing that generates and a mapping of string node ID to the fire probability value
func (meshping *MeshPing) GenerateProbabilitydatamap() map[string]float64 {
	// Create an empty probability map
	probdata := make(map[string]float64)

	// Iterate over the Pings of the meshping
	for nodeid, sensorping := range meshping.Pings {
		// Convert the nodeIDs to strings and assign the Fireprobability value
		probdata[strconv.FormatInt(nodeid, 10)] = float64(sensorping.Fireprobability)
	}

	// Return the probdata
	return probdata
}

// A method of MeshPing that flushes a completed MeshPing to the cloud.
func (meshping *MeshPing) Flush(meshorchestrator *MeshOrchestrator) error {
	// Generate a new PingDocument from the meshping
	pingdoc := NewPingDocument(meshping)

	// Push the pingdoc to the cloud and check the success.
	err := pingdoc.Push(&meshorchestrator.Cloudinterface)
	if err != nil {
		// Log the meshping failing to be flushed to the cloud.
		logmessage := NewOrchCloudlog(fmt.Sprintf("(failure) mesh ping accumulated and flush failed | doc - %v", meshping.PingID))
		meshorchestrator.LogQueue <- logmessage
	}

	// Log the meshping succesfully being flushed to the cloud.
	logmessage := NewOrchCloudlog(fmt.Sprintf("(success) mesh ping accumulated and flush successful | doc - %v", meshping.PingID))
	meshorchestrator.LogQueue <- logmessage

	// Delete the meshping from the accumulation
	delete(meshorchestrator.Accumulation, meshping.PingID)
	return nil
}

// A method of MeshPing that assigns a SensorPing to the MeshPing
// and flushes the MeshPing to the cloud if it is completed.
func (meshping *MeshPing) AddPing(sensorping SensorPing, meshorchestrator *MeshOrchestrator) {
	// Assign the SensorPing to the MeshPing's Pings map
	meshping.Pings[sensorping.Sensornode.NodeID] = sensorping

	// Check if the meshping is complete
	if meshping.Complete() {
		// Flush the mesh ping to the cloud
		meshping.Flush(meshorchestrator)
	}
}

// A function that handles the output of the SensorPings recieved over the meshorchestrator's AccumulatorQueue.
// Assigns the recieved to ping to the appropriate MeshPing in the orchestrator's accumulation or creates a new
// Mesh and assigns it to that new MeshPing and adds the new MeshPing to orchestrator's accumulation.
func PingHandler(meshorchestrator *MeshOrchestrator) {
	// log the beginning of the pinghandler
	meshorchestrator.LogQueue <- NewOrchServerlog("(startup) ping handler has started")

	// Iterate over the AccumulatorQueue until it closes.
	for sensorping := range meshorchestrator.AccumulatorQueue {

		// Check if the sensorping's ping ID exists on the accumulation
		if meshping, ok := meshorchestrator.Accumulation[sensorping.PingID]; ok {
			// Assign the sensor ping to existing meshping
			meshping.AddPing(sensorping, meshorchestrator)

		} else {
			// Create a new mesh ping with the sensorping's ping ID and ping time.
			meshping := *NewMeshPing(sensorping.PingID, sensorping.Pingtime, meshorchestrator.NodeIDlist)
			// Assign the sensorping into the new meshping
			meshping.AddPing(sensorping, meshorchestrator)
			// Add the new meshping into the meshorchestrator's accumulation.
			meshorchestrator.Accumulation[sensorping.PingID] = meshping
		}
	}
}
