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
			fmt.Printf("[error] config file could not be read - %v\n", err)
			fmt.Println("[suggestion] run 'fyrcli config -m generate' if file does not exist or is corrupted.")
			return
		}

		// Check the device type config value.
		if config.DeviceType != "mesh-controller" {
			fmt.Println("[error] server boot can only be performed on the mesh controller.")
			return
		}

		// Check if the 'FYRMESHSRC' env variable has been set.
		if srcdir == "" {
			fmt.Println("[error] server boot failed - environment variable 'FYRMESHSRC' has not set.")
			return
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
			fmt.Println("[error] unsupported server name -", server)
		}
	},
}

func bootORCH() {
	// Define the command to start the ORCH server in an lxterminal window.
	cmd := exec.Command("lxterminal", "--geometry=250x30", "-t", "ORCH", "-e", "fyrorch")
	// Run the command.
	cmd.Run()
}

func bootLINK(srcdir string) {
	// Define the path to the LINK server python script.
	linkserver := fmt.Sprintf("%v/fyrlink/interface.py", srcdir)
	// Define the command to start the LINK server python script.
	command := fmt.Sprintf("python3 %v", linkserver)
	// Define the command to start the LINK server in an lxterminal window.
	cmd := exec.Command("lxterminal", "--geometry=20x10", "-t", "LINK", "-e", command)
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
