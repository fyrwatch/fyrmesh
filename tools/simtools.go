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
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// A struct that represents a simulator
// seed for a given type of sensor value
type SimulatorSeed struct {
	// A float64 that  represents the initial value of the seed
	Initial float64
	// A float64 that represents the value by which the seed moves the cursor
	Adjust float64
	// A float64 that represents the width of the seed while generating values
	Width float64
	// A float64 that represents the max deviation from the initial value before the seed curve ends
	Peak float64
	// A float64 that represents the current value of the seed
	Cursor float64
	// A string that represents the kind of seed curve to follow
	Curve string
}

// A constructor function that generates and returns a SimulatorSeed object.
// Requires the initial, peak, adjust, width and curve name.
// The cursor is set the value of the initial passed.
func NewSimulatorSeed(initial, peak, adjust, width float64, curve string) *SimulatorSeed {
	// Create a new SimulatorSeed
	seed := SimulatorSeed{}

	// Assign the values
	seed.Initial = initial
	seed.Peak = peak
	seed.Adjust = adjust
	seed.Width = width
	seed.Cursor = initial
	seed.Curve = curve

	// Return the seed
	return &seed
}

// A method of SimulatorSeed that serves as the seed curve while rising.
// Increments the Cursor by Adjust until it exceeds the peak.
// However if the reverse bool is set, the opposite occurs.
func (seed *SimulatorSeed) rising(reverse bool) {
	for {
		// Check the cursor state
		if seed.Cursor > seed.Peak && !reverse {
			break
		} else if seed.Cursor < seed.Peak && reverse {
			break
		}

		// Increment the cursor
		seed.Cursor = seed.Cursor + seed.Adjust
		time.Sleep(time.Second * 5)
	}
}

// A method of SimulatorSeed that serves as the seed curve while falling.
// Decrements the Cursor by Adjust until it goes below the initial.
// However if the reverse bool is set, the opposite occurs.
func (seed *SimulatorSeed) falling(reverse bool) {
	for {
		// Check the cursor state
		if seed.Cursor < seed.Initial && !reverse {
			break
		} else if seed.Cursor > seed.Initial && reverse {
			break
		}

		// Decrement the cursor
		seed.Cursor = seed.Cursor - seed.Adjust
		time.Sleep(time.Second * 5)
	}
}

// A method of SimulatorSeed that serves as a flip curve.
// Simply flips the Cursor from 1 to 0 after the Adjust*5 amount of seconds
// and then flips it back after the same amount of time.
func (seed *SimulatorSeed) flip() {
	time.Sleep(time.Second * time.Duration(seed.Adjust*5))
	seed.Cursor = 1
	time.Sleep(time.Second * time.Duration(seed.Adjust*5))
	seed.Cursor = 0
}

// A method of SimulatorSeed that starts the curve of the seed.
func (seed *SimulatorSeed) StartCurve(wg *sync.WaitGroup) {
	// Check the type of curve to start
	switch seed.Curve {
	case "bell":
		seed.rising(false)
		time.Sleep(time.Second * 10)
		seed.falling(false)

	case "revbell":
		seed.rising(true)
		time.Sleep(time.Second * 10)
		seed.falling(true)

	case "flip":
		seed.flip()
	}

	// Decrement the waitgroup.
	wg.Done()
}

// A struct that represents a Fire Event Simulator
type FireEventSimulator struct {
	// A bool tht represents if the simulator is on.
	SimulationOn bool
	// A pool of SimulatorSeeds for each sensor type.
	SimulationSeeds map[string]*SimulatorSeed
}

// A constructor function that generates and returns a FireEventSimulator
// Creates the default seeds for each sensor type.
func NewFireEventSimulator() *FireEventSimulator {
	// Create a FireEventSimulator
	simulator := FireEventSimulator{}
	// Set the simulator to off
	simulator.SimulationOn = false
	// Create an empty map and assign it
	simulator.SimulationSeeds = make(map[string]*SimulatorSeed)

	// Create the SimulatorSeeds for each sensor type.
	simulator.SimulationSeeds["GAS"] = NewSimulatorSeed(450.0, 900.0, 25.0, 75.0, "bell")
	simulator.SimulationSeeds["HUM"] = NewSimulatorSeed(50.0, 22.5, -1.5, 3.0, "revbell")
	simulator.SimulationSeeds["TEM"] = NewSimulatorSeed(27.0, 55, 1.5, 3.0, "bell")
	simulator.SimulationSeeds["FLM"] = NewSimulatorSeed(0, 1, 12.0, 0, "flip")

	return &simulator
}

// A method of FireEventSimulator that starts a Fire Event.
// Requires a LogQueue to log the start and end of the Fire Event.
// Uses a wait group to monitor the completion of each individual seed's event curve.
func (simulator *FireEventSimulator) StartFireEvent(logqueue chan Log) {
	// Create a wait group
	wg := sync.WaitGroup{}
	// Log the start of the fire event
	logqueue <- NewOrchSchedlog("(simulator) fire event has started")

	// Iterate over the Seed pool
	for _, seed := range simulator.SimulationSeeds {
		// Increment the wait group
		wg.Add(1)
		// start the seed curve
		go seed.StartCurve(&wg)
	}

	// Wait for wait group to complete
	wg.Wait()
	// Log the end of the fire event
	logqueue <- NewOrchSchedlog("(simulator) fire event has ended")
}

// A method of FireEventSimulator that returns a
// simulated value for a given sensor type.
func (simulator *FireEventSimulator) GetSimulatedValue(sensortype string) float64 {
	// Create a float
	var simvalue float64

	// Check the sensor type and retrieve the appropriate seed.
	// Use that seed and its values to generate a random simulated value.
	switch sensortype {
	case "HUM":
		seed := simulator.SimulationSeeds["HUM"]
		simvalue = generaterandomvalue(seed.Cursor, seed.Cursor-seed.Width, 2)

	case "TEM":
		seed := simulator.SimulationSeeds["TEM"]
		simvalue = generaterandomvalue(seed.Cursor, seed.Cursor+seed.Width, 2)

	case "GAS":
		seed := simulator.SimulationSeeds["GAS"]
		simvalue = generaterandomvalue(seed.Cursor, seed.Cursor+seed.Width, 0)

	case "FLM":
		seed := simulator.SimulationSeeds["FLM"]
		simvalue = seed.Cursor
	}

	// Return the simulated value
	return simvalue
}

// A function that generates a sensor value given the the value as an unparsed string and the sensor type.
// The mesh orchestrator is used to retrieve data from the simulator when it is on.
// The enableCorrection flag allows the sensor values to be overriden when the sensor is defunct or for debugging.
func GenerateSensorValue(sensorvalue string, sensortype string, meshorchestrator *MeshOrchestrator, enableCorrection bool) float64 {
	// Create a float
	var generatedvalue float64

	// Check if simulation is on or if corrections have been enabled
	if meshorchestrator.Simulator.SimulationOn || enableCorrection {
		// generate a simulated value, either for a fire event or for baseline seed.
		generatedvalue = meshorchestrator.Simulator.GetSimulatedValue(sensortype)

	} else {
		// parse the provided string sensor value
		parsedval, _ := strconv.ParseFloat(sensorvalue, 64)
		generatedvalue = parsedval
	}

	// Return the generated value
	return generatedvalue
}

// A function that generates a random float64 number between a given two number with the precision set.
// The values are sorted for magnitude and the value is generated from the random seed.
// The precision must be a value between 0 and 8 and sets the number of decimal point in the float.
func generaterandomvalue(val1, val2 float64, precision int8) float64 {
	// Seed the random generator
	rand.Seed(time.Now().UnixNano())

	// Check the precision limit
	if precision > 8 || precision < 0 {
		return 0
	}

	// Declare the min and max
	var min, max float64
	// Determine the min and max
	if val1 > val2 {
		min = val2
		max = val1
	} else {
		min = val1
		max = val2
	}

	// Calculate a trunctuator for the given precision
	trunct := math.Pow(10, float64(precision))
	// Generate the random value and round it to the precision.
	randomvalue := min + rand.Float64()*(max-min)
	randomvalue = math.Round(randomvalue*trunct) / trunct
	// Return the generated value.
	return randomvalue
}
