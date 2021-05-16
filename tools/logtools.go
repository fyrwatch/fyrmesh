package tools

import (
	"fmt"
	"time"
)

// A function that returns the current time as an ISO8601 string without the timezone.
func currenttimeISO() string {
	return time.Now().Format("2006-01-02T15:04:05")
}

// A function that returns the parsed ORCH log string. Accepts a message,
// generates the time and constructs the log with the source always set to ORCH.
func GenerateORCHLog(message string) string {
	return fmt.Sprintf("[ORCH][%v] %v", currenttimeISO(), message)
}

// A function that returns the parsed LINK log string. Accepts a
// source, time and message and constructs the log accordingly.
func GenerateLINKLog(source string, logtime string, message string) string {
	return fmt.Sprintf("[%v][%v] %v", source, logtime, message)
}

// A function that handles the output of the logs recieved
// over a given logqueue. Currently only prints to stdout.
func LogHandler(logqueue chan string, observerqueue chan string) {
	// Declare the observer toggle
	observertoggle := false

	// Iterate over the logqueue until it closes.
	for log := range logqueue {

		// Check the logqueue for potential command codes
		switch log {
		case "enable-observe":
			// Enable the observerqueue
			observertoggle = true
			fmt.Printf("[LOG][%v] Observer Queue has been enabled.\n", currenttimeISO())

		case "disable-observe":
			// Disable the observerqueue
			observertoggle = false
			fmt.Printf("[LOG][%v] Observer Queue has been disabled.\n", currenttimeISO())

		default:
			// Regularly print the log to console.
			fmt.Println(log)

			// Check the toggle
			if observertoggle {
				// Pass every log into the observerqueue channel
				observerqueue <- log
			}
		}
	}
}
