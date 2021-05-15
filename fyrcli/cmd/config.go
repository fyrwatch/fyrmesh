/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	orch "github.com/fyrwatch/fyrmesh/fyrorch/orchpkg"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View/Manipulate configuration values of the FyrCLI.",
	Long: `View/Manipulate configuration values of the FyrCLI which are obtained from a configuration file.

The path to the configuration file is set in the environment variable 'FYRMESHCONFIG'. 
The command expects the 'mode' flag which defaults to 'show'. 

The list of valid manipulation modes are below:
- 'show'      - displays the configuration values from the file.
- 'locate'    - displays the path to the configuration file obtained from the 'FYRMESHCONFIG' env variable.
- 'checkfile' - confirms whether the configuration file currently exists.
- 'generate'  - generates a new configuration file with the default values (overwrite prompt will appear).
- 'modify'    - starts an interactive menu based shell to modify specific values of the configuration file.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the manipulation mode value from the command flags.
		mode, _ := cmd.Flags().GetString("mode")

		// Check the value of the manipulation mode name and call the appropriate method.
		switch mode {
		case "generate":
			configfunc_generate()
		case "checkfile":
			configfunc_checkfile()
		case "locate":
			configfunc_locate()
		case "show":
			configfunc_show()
		case "modify":
			configfunc_modify()
		default:
			fmt.Println("Unsupported config manipulation mode -", mode)
		}
	},
}

func configfunc_generate() {
	// Call the method to check if the config file exists.
	check, _ := orch.CheckConfig()
	if check {
		// File already exists. Request user confirmation to overwrite.
		fmt.Println("A config file already exists. Proceed to overwrite? y/n")
		// Read the user input.
		var proceed string
		fmt.Scanln(&proceed)

		// Test the vale of the user input.
		switch proceed {
		case "y", "Y", "yes", "Yes", "YES":
			// Overwrite approved.
			fmt.Println("Overwrite allowed. Existing config file will now be overwritten and reset to defaults.")

		case "n", "N", "no", "No", "NO":
			// Overwrite denied.
			fmt.Println("Overwrite cancelled.")
			return

		default:
			// Invalid Response
			fmt.Println("Invalid Response!")
			return
		}
	}

	// Call the method to generate a new default config file.
	err := orch.GenerateConfig()
	if err != nil {
		// Config failed to be generated.
		fmt.Printf("A config file could not be generated - %v\n", err)
	} else {
		// Config has been generated.
		fmt.Println("A config file has been generated.")
		// Print out some other suggested methods for the CLI tool.
		fmt.Println("-- Use 'fyrcli config -m show' to view the configuration values.")
		fmt.Println("-- Use 'fyrcli config -m locate' to view the path to the config file.")
	}
}

func configfunc_checkfile() {
	// Call the method to check if the config file exists.
	check, err := orch.CheckConfig()
	if check {
		// The file has been confirmed to exist
		fmt.Println("A configuration file exists")
		// Print out some other suggested methods for the CLI tool.
		fmt.Println("-- Use 'fyrcli config -m show' to view the configuration values.")
		fmt.Println("-- Use 'fyrcli config -m locate' to view the path to the config file.")
	} else {
		// The file either does not exist or there is some uncertainty in its existence.
		fmt.Println("A configuration file does not exist. It may also be corrupted or inaccesible.", err)
		// Print out some other suggested methods for the CLI tool.
		fmt.Println("-- Use 'fyrcli config -m generate' to generate a new configuration file.")
	}
}

func configfunc_locate() {
	// Retrieve the env variable 'FYRMESHCONFIG'.
	configpath := os.Getenv("FYRMESHCONFIG")
	// Print the value of the env variable.
	fmt.Println(configpath)
	// Print out some other suggested methods for the CLI tool.
	fmt.Println("-- Use 'fyrcli config -m show' to view the configuration values.")
}

func configfunc_show() {
	// Read the config file.
	config, err := orch.ReadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// Print out the configuration values in a formatted menu-type list.
	fmt.Println()
	fmt.Println("---- FyrMesh Configuration File ----")
	fmt.Println()

	fmt.Println("-- General Configuration --")
	fmt.Printf("Device ID: %v\n", config.DeviceID)
	fmt.Printf("Device Type: %v\n", config.DeviceType)
	fmt.Println()

	fmt.Println("-- ORCH Configuration --")
	fmt.Printf("Host: %v\n", config.Services["ORCH"].Host)
	fmt.Printf("Port: %v\n", config.Services["ORCH"].Port)
	fmt.Println()

	fmt.Println("-- LINK Configuration --")
	fmt.Printf("Host: %v\n", config.Services["LINK"].Host)
	fmt.Printf("Port: %v\n", config.Services["LINK"].Port)
	fmt.Println()

	fmt.Println("---- end of file ----")
	fmt.Println()

	// Print out some other suggested methods for the CLI tool.
	fmt.Println("-- Use 'fyrcli config -m modify' to modify the configuration values.")
	fmt.Println("-- Use 'fyrcli config -m locate' to view the path to the config file.")
}

func configfunc_modify() {
	currentconfig, err := orch.ReadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println()
	fmt.Println("---- FyrMesh Configuration Modifier ----")
	fmt.Println()

	fmt.Println("--------------------------------------------------------------")
	fmt.Println("Choose the type of value that needs to be modified. Enter the number")
	fmt.Println("1. General Configuration Values")
	fmt.Println("2. ORCH Configuration Values")
	fmt.Println("3. LINK Configuration Values")
	fmt.Println("--------------------------------------------------------------")

	var menunumber int
	newconfig := currentconfig
	fmt.Scanln(&menunumber)

	switch menunumber {
	case 1:
		fmt.Println("--------------------------------------------------------------")
		fmt.Println("Choose the General Configuration Value that needs to be changed")
		fmt.Println("1. Device ID")
		fmt.Println("2. Device Type")
		fmt.Println("--------------------------------------------------------------")
		fmt.Scanln(&menunumber)

		switch menunumber {
		case 1:
			fmt.Println("The Device ID is a fixed hardware value and cannot be changed")
			os.Exit(0)
		case 2:
			fmt.Println("The Device Type is a fixed hardware value and cannot be changed")
			os.Exit(0)
		default:
			fmt.Println("Invalid Choice. Start over!")
			os.Exit(0)
		}

	case 2:
		fmt.Println("--------------------------------------------------------------")
		fmt.Println("Choose the ORCH Configuration Value that needs to be changed")
		fmt.Println("1. Host URL")
		fmt.Println("2. Host Port")
		fmt.Println("--------------------------------------------------------------")
		fmt.Scanln(&menunumber)

		switch menunumber {
		case 1:
			var hosturl string
			fmt.Printf("The current value of ORCH Host URL is '%v'. Enter the new value(0 to not make a change)\n", currentconfig.Services["ORCH"].Host)
			fmt.Scanln(&hosturl)

			if hosturl != "0" {
				newconfig.Services["ORCH"] = orch.ServiceConfig{Host: hosturl, Port: currentconfig.Services["ORCH"].Port}
				orch.WriteConfig(newconfig)
			}
			os.Exit(0)

		case 2:
			var hostport int
			fmt.Printf("The current value of ORCH Host Port is '%v'. Enter the new value(0 to not make a change)\n", currentconfig.Services["ORCH"].Port)
			fmt.Scanln(&hostport)

			if hostport != 0 {
				newconfig.Services["ORCH"] = orch.ServiceConfig{Host: currentconfig.Services["ORCH"].Host, Port: hostport}
				orch.WriteConfig(newconfig)
			}
			os.Exit(0)

		default:
			fmt.Println("Invalid Choice. Start over!")
			os.Exit(0)
		}

	case 3:
		fmt.Println("--------------------------------------------------------------")
		fmt.Println("Choose the LINK Configuration Value that needs to be changed")
		fmt.Println("1. Host URL")
		fmt.Println("2. Host Port")
		fmt.Println("--------------------------------------------------------------")
		fmt.Scanln(&menunumber)

		switch menunumber {
		case 1:
			var hosturl string
			fmt.Printf("The current value of LINK Host URL is '%v'. Enter the new value(0 to not make a change)\n", currentconfig.Services["LINK"].Host)
			fmt.Scanln(&hosturl)

			if hosturl != "0" {
				newconfig.Services["LINK"] = orch.ServiceConfig{Host: hosturl, Port: currentconfig.Services["LINK"].Port}
				orch.WriteConfig(newconfig)
			}
			os.Exit(0)

		case 2:
			var hostport int
			fmt.Printf("The current value of LINK Host Port is '%v'. Enter the new value(0 to not make a change)\n", currentconfig.Services["LINK"].Port)
			fmt.Scanln(&hostport)

			if hostport != 0 {
				newconfig.Services["LINK"] = orch.ServiceConfig{Host: currentconfig.Services["LINK"].Host, Port: hostport}
				orch.WriteConfig(newconfig)
			}
			os.Exit(0)

		default:
			fmt.Println("Invalid Choice. Start over!")
			os.Exit(0)
		}
	default:
		fmt.Println("Invalid Choice. Start over!")
		os.Exit(0)
	}
}

func init() {
	// Add the command 'config' to root CLI command.
	rootCmd.AddCommand(configCmd)

	// Add the flag 'mode'.
	configCmd.Flags().StringP("mode", "m", "show", "mode of config manipulation.")
}
