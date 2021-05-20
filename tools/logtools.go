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
	"strings"
	"time"
)

// A function that returns the current time as an ISO8601 string without the timezone.
func CurrentISOtime() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05")
}

// A function that deserializes a a string with a format akin
// to 'key1-value1=key2-value2..' into a map[string]string """
func Deepdeserialize(str string) map[string]string {
	// Split the string into individual key-value pairs
	pairs := strings.Split(str, "=")
	// Define a new map[string]string oject
	dict := make(map[string]string)

	// Iterate over the key-value pairs
	for _, pair := range pairs {
		// Split each key-value pair
		set := strings.Split(pair, "-")
		// Add the key value pair into the map
		dict[set[0]] = set[1]
	}

	// Return the map
	return dict
}

// A struct that defines a log that is
// generated within the orchestrator.
type OrchLog struct {
	Logsource   string
	Logtype     string
	Logtime     string
	Logmessage  string
	Logmetadata map[string]string
}

// A struct that defines a log that is
// sent out of the server for observation.
// The Logmessage here is a fully stringified Log.
// The Logsource and Logtype are used for filtering.
type ObserverLog struct {
	Logsource  string
	Logtype    string
	Logmessage string
}

// An interface that defines a common interface that
// can be used by OrchLog and the ComplexLog proto
type Log interface {
	GetLogsource() string
	GetLogtype() string
	GetLogtime() string
	GetLogmessage() string
	GetLogmetadata() map[string]string
}

// A getter function for the Logsource field of OrchLog
func (orchlog *OrchLog) GetLogsource() string {
	if orchlog != nil {
		return orchlog.Logsource
	}
	return ""
}

// A getter function for the Logtype field of OrchLog
func (orchlog *OrchLog) GetLogtype() string {
	if orchlog != nil {
		return orchlog.Logtype
	}
	return ""
}

// A getter function for the Logtime field of OrchLog
func (orchlog *OrchLog) GetLogtime() string {
	if orchlog != nil {
		return orchlog.Logtime
	}
	return ""
}

// A getter function for the Logmessage field of OrchLog
func (orchlog *OrchLog) GetLogmessage() string {
	if orchlog != nil {
		return orchlog.Logmessage
	}
	return ""
}

// A getter function for the Logmetadata field of OrchLog
func (orchlog *OrchLog) GetLogmetadata() map[string]string {
	if orchlog != nil {
		return orchlog.Logmetadata
	}
	return nil
}

// A constructor function that generates and returns an OrchLog with
// the 'serverlog' type. The message passed is set as the Logmessage.
func NewOrchServerlog(message string) *OrchLog {
	// Construct a new OrchLog
	orchlog := OrchLog{}
	// Set the values of the OrchLog
	orchlog.Logsource = "ORCH"
	orchlog.Logtype = "serverlog"
	orchlog.Logtime = CurrentISOtime()
	orchlog.Logmessage = message
	orchlog.Logmetadata = make(map[string]string)
	// Return the OrchLog
	return &orchlog
}

// A constructor function that generates and returns an OrchLog with the
// 'protolog' type. The message passed is set as the Logmessage and the
// server, service and err values are set in the Logmetadata map.
func NewOrchProtolog(message string, server string, service string, err error) *OrchLog {
	// Construct a new OrchLog
	orchlog := OrchLog{}
	// Set the values of the OrchLog
	orchlog.Logsource = "ORCH"
	orchlog.Logtype = "protolog"
	orchlog.Logtime = CurrentISOtime()
	orchlog.Logmessage = message
	orchlog.Logmetadata = make(map[string]string)
	// Set the values of the OrchLog Metadata
	orchlog.Logmetadata["server"] = server
	orchlog.Logmetadata["service"] = service
	orchlog.Logmetadata["error"] = fmt.Sprintf("%v", err)
	// Return the OrchLog
	return &orchlog
}

// A constructor function that generates and returns an OrchLog with
// the 'cloudlog' type. The message passed is set as the Logmessage.
func NewOrchCloudlog(message string) *OrchLog {
	// Construct a new OrchLog
	orchlog := OrchLog{}
	// Set the values of the OrchLog
	orchlog.Logsource = "ORCH"
	orchlog.Logtype = "cloudlog"
	orchlog.Logtime = CurrentISOtime()
	orchlog.Logmessage = message
	orchlog.Logmetadata = make(map[string]string)
	// Return the OrchLog
	return &orchlog
}

// A constructor function that generates and returns an OrchLog with
// the 'schedlog' type. The message passed is set as the Logmessage.
func NewOrchSchedlog(message string) *OrchLog {
	// Construct a new OrchLog
	orchlog := OrchLog{}
	// Set the values of the OrchLog
	orchlog.Logsource = "ORCH"
	orchlog.Logtype = "schedlog"
	orchlog.Logtime = CurrentISOtime()
	orchlog.Logmessage = message
	orchlog.Logmetadata = make(map[string]string)
	// Return the OrchLog
	return &orchlog
}

// A constructor function that generates and returns an OrchLog with
// the 'obstoggle' type. The command passed is set as the Logmessage.
func NewObsCommand(command string) *OrchLog {
	// Construct a new OrchLog
	orchlog := OrchLog{}
	// Set the values of the OrchLog
	orchlog.Logsource = "OBS"
	orchlog.Logtype = "toggle"
	orchlog.Logtime = CurrentISOtime()
	orchlog.Logmessage = command
	orchlog.Logmetadata = make(map[string]string)
	// Return the OrchLog
	return &orchlog
}

// A constructore function that generates and returns an ObserverLog that is
// built using an existing Log. The Logmessage is a stringified version of the
// Log object being used. The Logsource and Logtype are taken from the Log struct.
func NewObserverLog(log Log) *ObserverLog {
	// Construct a new ObserverLog
	obslog := ObserverLog{}
	// Set the value of Logsource and Logtype
	obslog.Logsource = log.GetLogsource()
	obslog.Logtype = log.GetLogtype()
	// Stringify and set the Logmessage
	obslog.Logmessage = FormatLog(log)
	// Return the ObserverLog
	return &obslog
}

// A function that simplifies and formats a Log into a string.
// Every logtype has a different format but the general structure
// of the string log is - '[source][time][type] message. metadata..'
func FormatLog(log Log) string {
	// Retrieve all the Log data
	logsource := log.GetLogsource()
	logtype := log.GetLogtype()
	logtime := log.GetLogtime()
	logmessage := log.GetLogmessage()
	logmetadata := log.GetLogmetadata()

	// Declare a string log
	var strlog string
	// Define the common prefix of all logs
	logprefix := fmt.Sprintf("[%s][%s]%11s", logsource, logtime, logtype)

	// Check the logtype and set the appropriate format
	switch logtype {
	case "serverlog", "cloudlog", "schedlog":
		strlog = fmt.Sprintf("%v || %v |", logprefix, logmessage)

	case "protolog":
		strlog = fmt.Sprintf("%v || %v | server - %v | service - %v| error - %v |", logprefix, logmessage, logmetadata["server"], logmetadata["service"], logmetadata["error"])

	case "message":
		strlog = fmt.Sprintf("%v || (%v) %v | type - %v |", logprefix, logmetadata["format"], logmessage, logmetadata["type"])

	case "meshsync":
		strlog = fmt.Sprintf("%v || (meshevent) %v | event - %v |", logprefix, logmessage, logmetadata["synctype"])

	case "nodesync":
		strlog = fmt.Sprintf("%v || (meshevent) %v | offset - %v |", logprefix, logmessage, logmetadata["offset"])

	case "handshake":
		strlog = fmt.Sprintf("%v || (meshevent) %v | node - %v |", logprefix, logmessage, logmetadata["node"])

	case "sensordata":
		sensordata := Deepdeserialize(logmetadata["sensors"])
		strlog = fmt.Sprintf("%v || (data) %v | sensors - %v | node - %v | ping - %v |", logprefix, logmessage, sensordata, logmetadata["node"], logmetadata["ping"])

	case "configdata":
		strlog = fmt.Sprintf("%v || (data) %v | node - %v | ping - %v |", logprefix, logmessage, logmetadata["node"], logmetadata["ping"])

	case "ctrldata":
		strlog = fmt.Sprintf("%v || (data) %v | node - %v |", logprefix, logmessage, logmetadata["nodeID"])

	case "nodelist":
		strlog = fmt.Sprintf("%v || (data) %v", logprefix, logmessage)

	case "toggle":
		strlog = fmt.Sprintf("%v || (toggle) %v", logprefix, logmessage)
	}

	// Return the string log
	return strlog
}

// A function that handles the output of the logs recieved
// over a given logqueue. Currently only prints to stdout.
func LogHandler(meshorchestrator *MeshOrchestrator) {
	// Declare the observer toggle
	observertoggle := false
	// log the beginning of the loghandler
	fmt.Println(FormatLog(NewOrchServerlog("(startup) log handler has started")))

	// Iterate over the logqueue until it closes.
	for log := range meshorchestrator.LogQueue {

		// Check the source of the log
		logtype := log.GetLogtype()
		switch logtype {
		case "serverlog", "protolog", "cloudlog", "schedlog", "message", "nodesync":
			// Stringify and print
			fmt.Println(FormatLog(log))
			// Send into observer queue if toggle is set
			if observertoggle {
				meshorchestrator.ObserverQueue <- *NewObserverLog(log)
			}

		case "handshake", "meshsync":
			// Call the method to update the meshorchestrator's NodeIDlist
			go meshorchestrator.UpdateNodeIDlist()

			// Stringify and print
			fmt.Println(FormatLog(log))
			// Send into observer queue if toggle is set
			if observertoggle {
				meshorchestrator.ObserverQueue <- *NewObserverLog(log)
			}

		case "sensordata":
			// Set the sensor node data to be added into the accumulation queue
			go meshorchestrator.SetSensorData(log)

			// Stringify and print
			fmt.Println(FormatLog(log))
			// Send into observer queue if toggle is set
			if observertoggle {
				meshorchestrator.ObserverQueue <- *NewObserverLog(log)
			}

		case "configdata":
			// Set the node configuration on the meshorchestrator's Nodelist
			go meshorchestrator.SetNode(log)

			// Stringify and print
			fmt.Println(FormatLog(log))
			// Send into observer queue if toggle is set
			if observertoggle {
				meshorchestrator.ObserverQueue <- *NewObserverLog(log)
			}

		case "ctrldata":
			// Set the meshorchestrator's Controlnode
			go meshorchestrator.SetControlnode(log)

			// Stringify and print
			fmt.Println(FormatLog(log))
			// Send into observer queue if toggle is set
			if observertoggle {
				meshorchestrator.ObserverQueue <- *NewObserverLog(log)
			}

		case "nodelist":
			// Set the meshorchestrator's NodeIDlist
			go meshorchestrator.SetNodeIDlist(log)

			// Stringify and print
			fmt.Println(FormatLog(log))
			// Send into observer queue if toggle is set
			if observertoggle {
				meshorchestrator.ObserverQueue <- *NewObserverLog(log)
			}

		case "toggle":
			// Check the toggle command
			observertogglecommand := log.GetLogmessage()
			switch observertogglecommand {
			case "enable-observe":
				// Enable the observerqueue
				observertoggle = true
				// Generate a server log
				meshorchestrator.LogQueue <- NewOrchServerlog("observer queue enabled")

			case "disable-observe":
				// Disable the observerqueue
				observertoggle = false
				// Generate a server log
				meshorchestrator.LogQueue <- NewOrchServerlog("observer queue disabled")
			}
		}
	}
}
