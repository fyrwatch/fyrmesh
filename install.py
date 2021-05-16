"""
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.
 
This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed 
without the express permission of Manish Meganathan and Mariyam A.Ghani
===========================================================================
FyrMesh Installer Script

A python script that automates the installation of FyrMesh services.
===========================================================================
"""

import os

# Print console messages
print("[INFO] FyrMesh installation has begun.")
currentdir = os.path.dirname(os.path.abspath(__file__))

# FyrCLI install
clidir = os.path.join(currentdir, 'fyrcli')
os.chdir(clidir)
os.system('go install')
print("[INFO] FyrCLI installation done.")

# FyrORCH install
orchdir = os.path.join(currentdir, 'fyrorch')
os.chdir(orchdir)
os.system('go install')
print("[INFO] FyrORCH installation done.")

# Print console messages
print("[INFO] FyrLINK installation done.")
print("[INFO] FyrMesh installation completed.")
print()
print("[WARNING] The FYRMESHCONFIG and FYRMESHSRC environment variables must be set for services to operate.")
print("FYRMESHCONFIG - path to the FyrMesh configuration file.")
print("FYRMESHSRC - path to directory of FyrMesh source code.")
print("-- use 'fyrcli help' for the usage of the FyrCLI.")
print("-- use 'fyrcli boot -s LINK' to start the FyrLINK server.")
print("-- use 'fyrcli boot -s ORCH' to start the FyrORCH server.")
