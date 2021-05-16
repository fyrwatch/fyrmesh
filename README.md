# FyrMesh
![FyrMesh Banner](fyrmesh.png)
## A Go package for sensor mesh orchestration powered by gRPC and Protocol Buffers with a built-in FyrCLI. 

**Version: 0.2.0**  
**Platform: Raspbian OS and Windows**  
**Language: Go 1.16 and Python 3.8**

### **Contributors**
- **Manish Meganathan**
- **Mariyam A. Ghani**

### **Contents**
  - [**Overview**](#overview)
    - [**FyrLINK**](#fyrlink)
    - [**FyrORCH**](#fyrorch)
    - [**FyrCLI**](#fyrcli)
    - [``gopkg`` **fyrorch/orch**](#gopkg-fyrorchorch)
    - [``gopkg`` **tools**](#gopkg-tools)
    - [**gRPC and Protocol Buffers**](#grpc-and-protocol-buffers)
  - [**Installation**](#installation)
    - [**1. Prerequisites**](#1-prerequisites)
    - [**2. Install FyrMesh**](#2-install-fyrmesh)
    - [**3. Setup Environment Variables**](#3-setup-environment-variables)
      - [*Linux*](#linux)
      - [*Windows*](#windows)
  - [**Using the FyrCLI**](#using-the-fyrcli)
  - [**Starting the FyrMesh**](#starting-the-fyrmesh)
  - [**Related Material**](#related-material)

### **Overview**
This package is a collection of microservices and a CLI tool that collectively form the FyrMesh platform. The communication between microservices is done with protocol buffers over gRPC.
The various components of the package are described below.

#### **FyrLINK**  
A Python server that handles the communication between the Raspberry Pi and the Control Node
over the serial port. It exposes an interface with the ``LINK`` gRPC server which implements methods that allow clients to read from and write to the serial port of the Raspberry Pi while ensuring the integrity of the message structures expected by the control node and semi-parsing the messages recieved.

#### **FyrORCH**  
A Go server that handles the orchestration of the mesh and the communication to the cloud backend services and the database through Firebase. It exposes an interface with the ``ORCH`` gRPC server
which implements methods for various mesh functionality that manipulate the state of the mesh or send messages on it.

#### **FyrCLI**  
A command-line interface application written in Go with the [**Cobra**](https://github.com/spf13/cobra) framework. It contains commands that allow a user to interact with the mesh through the orchestrator over a gRPC connection. The CLI tool can be used on the mesh controller itself or even on a remote system within the local network that is configured as a mesh observer device.

#### ``gopkg`` **fyrorch/orch**   
A Go package that contains the implementation of the ``LINK`` gRPC client, the ``ORCH`` gRPC server and the ``ORCH`` gRPC client. The client implementations contain helper functions to call its service methods appropriately.

#### ``gopkg`` **tools**  
A Go package that contains tools used by various **FyrMesh** services such interacting with the config file, handling logging I/O, interfacing with cloud services and manipulating internal data structures.

#### gRPC and Protocol Buffers

### **Installation**
The FyrMesh library can be installed onto a Linux or a Windows system. Linux support is only intended for the Raspbian OS that is designed for the Raspberry Pi.  

#### 1. Prerequisites
Install Go 1.16 and Python 3.8
   - [Install Go on the Raspberry Pi](https://www.jeremymorgan.com/tutorials/raspberry-pi/install-go-raspberry-pi/)
   - [Install Go on Windows](https://golang.org/doc/install)
   - [Install Python on the Raspberry Pi](https://projects.raspberrypi.org/en/projects/generic-python-install-python3#:~:text=Open%20your%20web%20browser%20and,a%20download%20will%20start%20automatically.)
   - [Install Python on Windows](https://realpython.com/installing-python/)


Download this repository to some location on your system. This location will be used to configure ``FYRMESHSRC`` in the next steps.
   - Download the latest release from the [releases](https://github.com/fyrwatch/fyrmesh/releases).
   - Use the Git bash to clone the repository.  
  ``` 
  git clone https://github.com/fyrwatch/fyrmesh.git 
  ```
   - Use ``go get`` to download the package.
  ```
  go get -u github.com/fyrwatch/fyrmesh
  ```

#### 2. Install FyrMesh
- Navigate into the /fyrmesh directory of the repository after downloading it.
- Open a terminal window in this directory and run the following command
```
python3 install.py
```
- The FyrMesh components will now be installed.

#### 3. Setup Environment Variables 
The FyrMesh services and applications require the environment variables ``FYRMESHCONFIG`` and ``FYRMESHSRC``. These environment variables have to be persistent and define the path to the the config file and the path to fyrmesh source folder respectively. These variables while not required at installation are required for the application to work.

The steps to add environment variables are detailed below for each platform.

##### Linux
- Run the following to cd into the ``/etc/profile.d`` directory and create a new file ``fyrmesh.sh``. This directory contains the shell scripts that are run at startup.
  ```
  $ cd /etc/profile.d
  $ nano envvars.sh
  ``` 
- Add the following lines to ``fyrmesh.sh`` file.
  ```
  export FYRMESHCONFIG=<<path to the config file>>
  export FYRMESHSRC=<<path to the repository source files directory>>
  export PATH=$PATH:<<path to Go bin install directory>>
  ``` 
  The recommended value for ``FYRMESHCONFIG`` is ``/home/pi/.fyrmesh/config.json``  
  The recommended value for go bin path extension is ``/home/pi/go/bin``
- Reboot the system to effect changes.

##### Windows
- Press the Start Menu and type 'env'. 
- An option 'Edit the environment variables for your account' will appear. Select it.
- A list of environment variables will be visible. Select the 'New' option.
- Fill in the value for ``FYRMESHCONFIG``. The recommended value is ``C:\Users\<user>\.fyrmesh\config.json``.
- Repeat and fill in the value for ``FYRMESHSRC`` with the path to the install directory.
- PATH extension wouldn't be required because the Go installer for windows.
- Reboot the systems to effect changes.

The installation of FyrMesh is now complete.  
Run the ``fyrcli`` command on the console to try it.

### Using the FyrCLI

### Starting the FyrMesh

### Related Material
