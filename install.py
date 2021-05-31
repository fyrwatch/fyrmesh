"""
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
