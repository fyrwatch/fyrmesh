"""
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.
 
This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed 
without the express permission of Manish Meganathan and Mariyam A.Ghani
===========================================================================
FyrMesh Code Generator Script

A python script that automates the code generation of the Protocol 
Buffers and gRPC service libraries for the fyrmesh.proto file.
===========================================================================
"""

import os
from grpc_tools import protoc

# Generate code for Python
protoc.main(('', '-I.', '--python_out=.', '--grpc_python_out=.', 'proto/fyrmesh.proto',))

# Generate code for Golang
os.system('protoc --go_out . --go-grpc_out . proto/fyrmesh.proto')
