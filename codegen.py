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
