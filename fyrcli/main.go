/*
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.

This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed
without the express permission of Manish Meganathan and Mariyam A.Ghani
===========================================================================
FyrMesh FyrCLI
===========================================================================
*/
package main

import (
	"fmt"

	"github.com/fyrwatch/fyrmesh/fyrcli/cmd"
)

func main() {
	// Colors are defined with ANSI escape commands https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
	// Change terminal color to orange
	fmt.Println("\033[38;5;208m")

	cmd.Execute()

	// Change terminal color to original
	fmt.Println("\033[0m")
}
