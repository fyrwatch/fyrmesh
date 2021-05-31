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
FyrMesh gopkg tools
===========================================================================
*/
package tools

import (
	"fmt"
	"math"
	"strconv"
)

// A function that maps a given input range of numbers to an output range.
// Takes a value within the input range and returns the mapped output.
//
// Implementation logic based on https://stackoverflow.com/a/5732390
func rangemapper(value, ipstart, ipend, opstart, opend float64) float64 {
	// Checks if value is within input range
	if value < ipstart || value > ipend {
		return 0
	}
	// Calculates the mapping and returns it
	output := opstart + ((opend-opstart)/(ipend-ipstart))*(value-ipstart)
	return output
}

// A function that calculates the probability of a
// fire based on the temperature reading around it.
func calcTEMprobability(sensorvalue float64) float64 {
	// Declare a float
	var probability float64

	// Check range of the temperature value and apply an appropriate mapping
	if sensorvalue < 0 {
		probability = 0

	} else if sensorvalue >= 0 && sensorvalue <= 20 {
		probability = rangemapper(sensorvalue, 0, 20, 1, 10)

	} else if sensorvalue > 20 && sensorvalue <= 35 {
		probability = rangemapper(sensorvalue, 20, 35, 10, 50)

	} else if sensorvalue > 35 && sensorvalue <= 40 {
		probability = rangemapper(sensorvalue, 35, 40, 50, 80)

	} else if sensorvalue > 40 && sensorvalue <= 45 {
		probability = rangemapper(sensorvalue, 40, 45, 80, 90)

	} else if sensorvalue > 45 && sensorvalue <= 55 {
		probability = rangemapper(sensorvalue, 45, 55, 90, 99)

	} else {
		probability = 100
	}

	// Return the probability
	return probability
}

// A function that calculates the probability of a
// fire based on the humidity reading around it.
func calcHUMprobability(sensorvalue float64) float64 {
	// Declare a float
	var probability float64

	// Check range of the humidity percentage value and apply an appropriate mapping
	if sensorvalue > 90 {
		probability = 0

	} else if sensorvalue < 90 && sensorvalue >= 50 {
		probability = rangemapper(sensorvalue, 90, 50, 1, 30)

	} else if sensorvalue < 50 && sensorvalue >= 40 {
		probability = rangemapper(sensorvalue, 50, 40, 30, 50)

	} else if sensorvalue < 40 && sensorvalue >= 30 {
		probability = rangemapper(sensorvalue, 40, 30, 50, 70)

	} else if sensorvalue < 30 && sensorvalue >= 20 {
		probability = rangemapper(sensorvalue, 30, 20, 70, 85)

	} else if sensorvalue < 20 && sensorvalue >= 10 {
		probability = rangemapper(sensorvalue, 20, 10, 85, 99)

	} else {
		probability = 100
	}

	// Return the probability
	return probability
}

// A function that calculates the probability of a fire
// based on the concentration of smoke reading around it.
func calcGASprobability(sensorvalue float64) float64 {
	// Declare a float
	var probability float64

	// Check range of the gas concentration value and apply an appropriate mapping
	if sensorvalue < 250 {
		probability = 0

	} else if sensorvalue >= 250 && sensorvalue <= 450 {
		probability = rangemapper(sensorvalue, 250, 450, 1, 15)

	} else if sensorvalue > 450 && sensorvalue <= 650 {
		probability = rangemapper(sensorvalue, 450, 650, 15, 40)

	} else if sensorvalue > 650 && sensorvalue <= 800 {
		probability = rangemapper(sensorvalue, 650, 800, 40, 70)

	} else if sensorvalue > 800 && sensorvalue <= 900 {
		probability = rangemapper(sensorvalue, 800, 900, 70, 90)

	} else if sensorvalue > 900 && sensorvalue <= 1000 {
		probability = rangemapper(sensorvalue, 900, 1000, 90, 99)

	} else {
		probability = 100
	}

	// Return the probability
	return probability
}

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

// A function that calculates the probability of a fire in the neighbourhood
// of the node and sets it to the SensorPing object's Fireprobability field.
func (sensorping *SensorPing) CalculateFireProbability() error {
	// Create a float and an empty slice of float64
	var probability float64
	probabilityvalues := make([]float64, 0)

	// Iterate over the Sensordata map
	for sensortype, sensorvalue := range sensorping.Sensordata {

		// Check the type of sensor and call the function
		// to calculate the probability from the sensor value
		switch sensortype {
		case "TEM":
			probability = calcTEMprobability(sensorvalue)

		case "HUM":
			probability = calcHUMprobability(sensorvalue)

		case "GAS":
			probability = calcGASprobability(sensorvalue)

		case "FLM":
			if sensorvalue == 1 {
				probability = 100
			} else {
				probability = 0
			}
		}

		// Append the probability value for the given sensor into the slice
		probabilityvalues = append(probabilityvalues, probability)
	}

	// Accumulate the probability values into its sum
	total := 0.0
	for _, prob := range probabilityvalues {
		total = total + prob
	}

	// Calulate average of all sensor probabilities and round it the 2 decimals
	probability = total / float64(len(probabilityvalues))
	probability = math.Round(probability*100) / 100
	// Assign the probability value to the Fireprobability field.
	sensorping.Fireprobability = probability
	return nil
}

// A method of SensorPing that generates the sensor data from the map parsed from the sensordata log.
// Passes the sensor data through some regularity checks and through the simulator seed.
// This allows the simulator to override the values if necessary.
func (sensorping *SensorPing) GenerateSensorData(sensordata map[string]string, meshorchestrator *MeshOrchestrator) {
	// Create an empty map of string -> float64
	sensorping.Sensordata = make(map[string]float64)

	// Parse and generate the sensor values that exist
	sensorkeys := []string{"HUM", "TEM", "FLM", "GAS"}
	for _, sensortype := range sensorkeys {
		if sensorvalue, ok := sensordata[sensortype]; ok {
			genvalue := GenerateSensorValue(sensorvalue, sensortype, meshorchestrator, true)
			sensorping.Sensordata[sensortype] = genvalue
		}
	}

}

// A constructor function that generates and returns a SensorPing.
// Requires a Log of type 'sensordata'
// The value of Sensordata is collected by parsed the sensordata log.
// The value of the Sensornode is retrieved from MeshOrchestrator's Nodelist.
// The value of the PingID is taken from sensordata log.
// The value of the Pingtime is the current time when the struct is constructed.
// The value of Fireprobabiliy is set by the CalculateFireProbability method.
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
	// Generate the Sensordata by parsing the data recieved from the log
	sensorping.GenerateSensorData(sensordata, meshorchestrator)

	// Retrieve the node ID from the metadata
	nodeid, _ := strconv.ParseInt(metadata["node"], 0, 64)
	// Retrieve the SensorNode object for the nodeID from the mesh orchestrator
	sensorping.Sensornode = meshorchestrator.Nodelist[nodeid]
	// Assign the ping ID from the metadata
	sensorping.PingID = metadata["ping"]
	// Assign the ping time to the current time
	sensorping.Pingtime = CurrentISOtime()

	// Calculate the value of the fire probability
	sensorping.CalculateFireProbability()
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

// A method of MeshPing that generates the average probability value of the meshping
// by accumulating the probability values of individual nodes and taking their average.
func (meshping *MeshPing) GenerateAvgProbability() float64 {
	// Create an empty slice
	probabilityvalues := make([]float64, 0)

	// Iterate over the Pings of the meshping
	for _, sensorping := range meshping.Pings {
		// Convert the nodeIDs to strings and assign the Fireprobability value
		probabilityvalues = append(probabilityvalues, sensorping.Fireprobability)
	}

	// Accumulate the probability values of each node
	total := 0.0
	for _, prob := range probabilityvalues {
		total = total + prob
	}

	// Calculate the average and round to 2 decimals
	avgprobability := total / float64(len(probabilityvalues))
	avgprobability = math.Round(avgprobability*100) / 100

	// Return the avgprobability
	return avgprobability
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
