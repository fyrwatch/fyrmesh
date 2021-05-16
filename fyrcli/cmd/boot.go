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
	"os/exec"

	"github.com/spf13/cobra"

	tools "github.com/fyrwatch/fyrmesh/tools"
)

// bootCmd represents the boot command
var bootCmd = &cobra.Command{
	Use:   "boot",
	Short: "Boots a FyrMesh gRPC server.",
	Long: `Boots a FyrMesh gRPC server depending on the server name flag set (mandatory).

Server booting can only performed on a device configured as a 'mesh-controller', which are 
typically devices that run Linux on the ARM architecture such as the Raspberry Pi 4B. 
The server booting is done by spawning a new 'lxterminal' window and starting the server.

The valid values of for the server name flag are below:
- values such as 'ORCH', 'orch' and 'orchestrator' -> boot the ORCH server. 
- values such as 'LINK', 'link' and 'interface' -> boot the LINK server.

NOTE: The 'FYRMESHSRC' and 'FYRMESHCONFIG' env variables must be set for boot systems to work.
NOTE: The LINK server should be booted up before the ORCH server to avoid an error.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the server name value from the command flags.
		server, _ := cmd.Flags().GetString("server")
		// Retrieve the env variable 'FYRMESHSRC'.
		srcdir := os.Getenv("FYRMESHSRC")

		// Read the config file.
		config, err := tools.ReadConfig()
		if err != nil {
			fmt.Printf("Config file could not be read - %v\n", err)
			fmt.Println("Run 'fyrcli config -m generate' if file does not exist or is corrupted.")
			os.Exit(0)
		}

		// Check the device type config value.
		if config.DeviceType != "mesh-controller" {
			fmt.Println("Server boot can only be performed on the mesh controller.")
			os.Exit(0)
		}

		// Check if the 'FYRMESHSRC' env variable has been set.
		if srcdir == "" {
			fmt.Println("Server boot failed - environment variable 'FYRMESHCONFIG' has not set.")
			os.Exit(0)
		}

		// Check the value of the server name and call the appropriate boot method.
		switch server {
		case "ORCH", "orch", "orchestrator":
			// Boot the ORCH server.
			bootORCH()

		case "LINK", "link", "interface":
			// Boot the LINK server
			bootLINK(srcdir)

		default:
			fmt.Println("Unsupported server name -", server)
		}
	},
}

func bootORCH() {
	// Define the command to start the ORCH server in an lxterminal window.
	cmd := exec.Command("lxterminal", "-e", "fyrorch")
	// Run the command.
	cmd.Run()
}

func bootLINK(srcdir string) {
	// Define the path to the LINK server python script.
	linkserver := fmt.Sprintf("%v/fyrlink/interface.py", srcdir)
	// Define the command to start the LINK server python script.
	command := fmt.Sprintf("python3 %v", linkserver)
	// Define the command to start the LINK server in an lxterminal window.
	cmd := exec.Command("lxterminal", "-e", command)
	// Run the command.
	cmd.Run()
}

func init() {
	// Add the command 'boot' to root CLI command.
	rootCmd.AddCommand(bootCmd)

	// Add the flag 'server' and mark it as a required flag.
	bootCmd.Flags().StringP("server", "s", "", "name of the server to boot up.")
	bootCmd.MarkFlagRequired("server")
}
