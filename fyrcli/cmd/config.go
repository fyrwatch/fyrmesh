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
FyrMesh FyrCLI
===========================================================================
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	tools "github.com/fyrwatch/fyrmesh/tools"
)

// configCmd represents the 'config' command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View configuration values of the FyrCLI.",
	Long: `View configuration values of the FyrCLI which are obtained from a configuration file.
The configuration file is in the directory defined by the environment variable 'FYRMESHCONFIG'. `,

	Run: func(cmd *cobra.Command, args []string) {
		showconfig()
	},
}

// configGenerateCmd represents the 'config generate' command
var configGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a new configuration file.",
	Long:  `Generates a new configuration file with the default values (overwrite prompt will appear if the file already exists).`,

	Run: func(cmd *cobra.Command, args []string) {
		// Call the method to check if the config file exists.
		check, _ := tools.CheckConfig()
		if check {
			// File already exists. Request user confirmation to overwrite.
			fmt.Print("[prompt] a config file already exists. proceed to overwrite? [y/n] > ")
			// Read the user input.
			var proceed string
			fmt.Scanln(&proceed)

			// Test the vale of the user input.
			switch proceed {
			case "y", "Y", "yes", "Yes", "YES":
				// Overwrite approved.
				fmt.Println("[info] overwrite allowed. existing config file will now be overwritten and reset to defaults.")

			case "n", "N", "no", "No", "NO":
				// Overwrite denied.
				fmt.Println("[end] overwrite cancelled.")
				return

			default:
				// Invalid Response
				fmt.Println("[error] invalid response!")
				return
			}
		}

		// Call the method to generate a new default config file.
		err := tools.GenerateConfig()
		if err != nil {
			// Config failed to be generated.
			fmt.Println("[failure] a config file could not be generated.")
			fmt.Printf("[error] %v", err)
		} else {
			// Config has been generated.
			fmt.Println("[success] a config file has been generated.")
			// Print out some other suggested methods for the CLI tool.
			fmt.Println("\n[suggestion] -- use 'fyrcli config show' to view the configuration values.")
			fmt.Println("[suggestion] -- use 'fyrcli config locate' to view the path to the config file.")
		}
	},
}

// configCheckfileCmd represents the 'config checkfile' command
var configCheckfileCmd = &cobra.Command{
	Use:   "checkfile",
	Short: "Confirms the existence of the configuration file.",
	Long:  `Confirms the existence of the configuration file in directory specfied by the 'FYRMESHCONFIG' env variable.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Call the method to check if the config file exists.
		check, err := tools.CheckConfig()
		if check {
			// The file has been confirmed to exist
			fmt.Println("[true] a configuration file exists")
			// Print out some other suggested methods for the CLI tool.
			fmt.Println("\n[suggestion] -- use 'fyrcli config show' to view the configuration values.")
			fmt.Println("[suggestion] -- use 'fyrcli config locate' to view the path to the config file.")
		} else {
			// The file either does not exist or there is some uncertainty in its existence.
			fmt.Println("[false] a configuration file does not exist. the file may also be corrupted or inaccesible.", err)
			// Print out some other suggested methods for the CLI tool.
			fmt.Println("\n[suggestion] -- use 'fyrcli config generate' to generate a new configuration file.")
		}
	},
}

// configCheckfileCmd represents the 'config checkfile' command
var configLocateCmd = &cobra.Command{
	Use:   "locate",
	Short: "Displays the path to the configuration file.",
	Long:  `Displays the path to the configuration file obtained from the 'FYRMESHCONFIG' env variable.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the env variable 'FYRMESHCONFIG'.
		configpath := os.Getenv("FYRMESHCONFIG")
		// Construct the path to the config file
		configfilepath := filepath.Join(configpath, "config.json")
		// Check if the env var has been set
		if configpath == "" {
			fmt.Println("[error] environment variable 'FYRMESHCONFIG' has not been set")
		}

		// Print the value of the env variable.
		fmt.Println(configfilepath)
		// Print out some other suggested methods for the CLI tool.
		fmt.Println("\n[suggestion] -- use 'fyrcli config show' to view the configuration values.")
	},
}

// configCheckfileCmd represents the 'config checkfile' command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays the value from the configuration file.",
	Long:  `Displays the value from the configuration file.`,

	Run: func(cmd *cobra.Command, args []string) {
		showconfig()
	},
}

// configCheckfileCmd represents the 'config checkfile' command
var configModifyCmd = &cobra.Command{
	Use:   "modify",
	Short: "Modify the specific values of the configuration file.",
	Long:  `Modify the specific values of the configuration file by starting an interactive shell.`,

	Run: func(cmd *cobra.Command, args []string) {
		currentconfig, err := tools.ReadConfig()
		if err != nil {
			fmt.Println(err)
			return
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
			fmt.Println("3. Scheduler Ping Rate")
			fmt.Println("--------------------------------------------------------------")
			fmt.Scanln(&menunumber)

			switch menunumber {
			case 1:
				fmt.Println("[info] the 'Device ID' is a fixed hardware value and cannot be changed")
				return
			case 2:
				fmt.Println("[info] the 'Device Type' is a fixed hardware value and cannot be changed")
				return
			case 3:
				var pingrate int
				fmt.Printf("[prompt] the current value of Scheduler Ping Rate is '%v'. Enter the new value (0 to not make a change)\n", currentconfig.SchedulerPingRate)
				fmt.Scanln(&pingrate)

				if pingrate != 0 {
					newconfig.SchedulerPingRate = pingrate
					tools.WriteConfig(newconfig)
				}

			default:
				fmt.Println("[error] invalid choice. start over!")
				return
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
				fmt.Printf("[prompt] the current value of 'ORCH Host URL' is '%v'. enter the new value (0 to not make a change)\n", currentconfig.Services["ORCH"].Host)
				fmt.Scanln(&hosturl)

				if hosturl != "0" {
					newconfig.Services["ORCH"] = tools.ServiceConfig{Host: hosturl, Port: currentconfig.Services["ORCH"].Port}
					tools.WriteConfig(newconfig)
				}
				return

			case 2:
				var hostport int
				fmt.Printf("[prompt] the current value of 'ORCH Host Port' is '%v'. Enter the new value (0 to not make a change)\n", currentconfig.Services["ORCH"].Port)
				fmt.Scanln(&hostport)

				if hostport != 0 {
					newconfig.Services["ORCH"] = tools.ServiceConfig{Host: currentconfig.Services["ORCH"].Host, Port: hostport}
					tools.WriteConfig(newconfig)
				}
				return

			default:
				fmt.Println("[error] invalid choice. start over!")
				return
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
				fmt.Printf("[prompt] the current value of 'LINK Host URL' is '%v'. enter the new value (0 to not make a change)\n", currentconfig.Services["LINK"].Host)
				fmt.Scanln(&hosturl)

				if hosturl != "0" {
					newconfig.Services["LINK"] = tools.ServiceConfig{Host: hosturl, Port: currentconfig.Services["LINK"].Port}
					tools.WriteConfig(newconfig)
				}
				return

			case 2:
				var hostport int
				fmt.Printf("[prompt] the current value of 'LINK Host Port' is '%v'. enter the new value (0 to not make a change)\n", currentconfig.Services["LINK"].Port)
				fmt.Scanln(&hostport)

				if hostport != 0 {
					newconfig.Services["LINK"] = tools.ServiceConfig{Host: currentconfig.Services["LINK"].Host, Port: hostport}
					tools.WriteConfig(newconfig)
				}
				return

			default:
				fmt.Println("[error] invalid choice. start over!")
				return
			}
		default:
			fmt.Println("[error] invalid choice. start over!")
			return
		}
	},
}

func showconfig() {
	// Read the config file.
	config, err := tools.ReadConfig()
	if err != nil {
		fmt.Printf("[error] %v", err)
		return
	}

	// Print out the configuration values in a formatted menu-type list.
	fmt.Println("---- FyrMesh Configuration File ----")
	fmt.Println()

	fmt.Println("-- General Configuration --")
	fmt.Printf("Device ID: %v\n", config.DeviceID)
	fmt.Printf("Device Type: %v\n", config.DeviceType)
	fmt.Printf("Scheduler Ping Rate: %v\n", config.SchedulerPingRate)
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

	// Print out some other suggested methods for the CLI tool.
	fmt.Println("\n[suggestion] -- use 'fyrcli config modify' to modify the configuration values.")
	fmt.Println("[suggestion] -- use 'fyrcli config locate' to view the path to the config file.")
}

func init() {
	// Add the command 'config' to root CLI command.
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configLocateCmd)
	configCmd.AddCommand(configGenerateCmd)
	configCmd.AddCommand(configCheckfileCmd)
	configCmd.AddCommand(configModifyCmd)
}
