package orch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// A struct that defines the configuration of the FyrMesh service
type Config struct {
	DeviceID   string                   `json:"deviceID"`
	DeviceType string                   `json:"deviceType"`
	Services   map[string]ServiceConfig `json:"services"`
}

// A struct that defines the configuration of an individual
// service that is a part of the FyrMesh service
type ServiceConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// A function that reads the config file that is located in the path specified
// in the 'FYRMESHCONFIG' env variable and returns the values in a Config struct.
func ReadConfig() (Config, error) {

	// Read the 'FYRMESHCONFIG' env var (no need to check if its set because the call to CheckConfig will handle that).
	filelocation := os.Getenv("FYRMESHCONFIG")

	// Check if the config file exists.
	check, err := CheckConfig()
	if !check {
		return Config{}, fmt.Errorf("config file does not exist - %v", err)
	}

	// Open the config file
	configFile, err := os.Open(filelocation)
	if err != nil {
		return Config{}, err
	}

	// Defer the closing of the file
	defer configFile.Close()

	// Read the config file into a byte array
	var config Config
	byteValue, _ := ioutil.ReadAll(configFile)

	// Marhsall the JSON byte array into a struct and return it
	json.Unmarshal([]byte(byteValue), &config)
	return config, nil
}

// A function that checks if the config file currently exists in the path specified
// in the 'FYRMESHCONFIG' env variable and returns the confirmation as a boolean.
func CheckConfig() (bool, error) {
	// Read the 'FYRMESHCONFIG' env var
	filelocation := os.Getenv("FYRMESHCONFIG")
	if filelocation == "" {
		return false, fmt.Errorf("environment variable 'FYRMESHCONFIG' has not set")
	}

	// Check if the file exists at the location
	if _, err := os.Stat(filelocation); err == nil {
		// File exists.
		return true, nil
	} else if os.IsNotExist(err) {
		// File does not exist.
		return false, nil
	} else {
		// File may or may not exist. (Schrodinger Case)
		return false, err
	}
}

// A function that writes a Config struct into a config file that is
// located in the path specified by the 'FYRMESHCONFIG' env variable.
func WriteConfig(config Config) error {
	// Read the 'FYRMESHCONFIG' env var
	filelocation := os.Getenv("FYRMESHCONFIG")
	if filelocation == "" {
		return fmt.Errorf("environment variable 'FYRMESHCONFIG' has not set")
	}

	// Format and Indent the config struct provided into a byte array.
	file, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return fmt.Errorf("could not format and marshal config struct - %v", err)
	}

	// Write the byte array to the filelocation.
	if err = ioutil.WriteFile(filelocation, file, 0644); err != nil {
		return fmt.Errorf("could not write config - %v", err)
	}

	return nil
}

// A function that extracts the Serial ID from the byte array
// output generated by the bash command to retrieve the CPU info
func extractserial(serialbytes []byte) string {
	// Parse the byte array into a string
	serial := fmt.Sprintf("%s", serialbytes)

	// Spilt the string by the : separator
	serialparts := strings.Split(serial, ":")

	// Extract the value element from the string
	serial = serialparts[1]

	// Trim the white space around the value and return it.
	serial = strings.TrimSpace(serial)
	return serial
}

// A function that generates the default configuration values and creates
// a new Config variable with those and writes this struct into a config file
// located in the path specified by the 'FYRMESHCONFIG' env variable.
func GenerateConfig() error {
	// Generate a default config with default values.
	defaultConfig := Config{
		DeviceID:   "unconfigured-device",
		DeviceType: "unconfigured-device",
		Services: map[string]ServiceConfig{
			"ORCH": {Host: "localhost", Port: 50001},
			"LINK": {Host: "localhost", Port: 50000},
		},
	}

	// Test the runtime environment and generate device values.
	if runtime.GOOS == "linux" && runtime.GOARCH == "arm" {
		// Assume the device is a Raspberry Pi system.

		// Read the unique Serial ID of the Raspberry Pi
		cmd := "cat /proc/cpuinfo | grep Serial"
		output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
		if err != nil {
			return fmt.Errorf("device serial could not be retrieved - %v", err)
		}

		// Extract the Serial ID from the bash output
		serial := extractserial(output)

		// Set the appropriate values to the default config
		defaultConfig.DeviceID = serial
		defaultConfig.DeviceType = "mesh-controller"

	} else if runtime.GOOS == "windows" {
		// The device is a remote observer Windows system.
		// Set the appropriate values to the default config
		defaultConfig.DeviceID = "observer-xxx"
		defaultConfig.DeviceType = "mesh-observer"

	} else {
		// FyrMesh is only meant to be operated in the
		// configurations above. Atleast for now..
		return fmt.Errorf("device runs on an unsupported OS and architecture")
	}

	// Write the generated config and check for errors.
	err := WriteConfig(defaultConfig)
	if err != nil {
		return fmt.Errorf("config write failed - %v", err)
	} else {
		return nil
	}
}
