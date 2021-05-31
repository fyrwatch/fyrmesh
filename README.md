# FyrMesh
![FyrMesh Banner](fyrmesh.png)
## A Go package for sensor mesh orchestration powered by gRPC & Protocol Buffers with a built-in FyrCLI. 

**Version: 0.3.0**  
**Platform: Raspbian OS & Windows**  
**Language: Go 1.16 & Python 3.9**
**License: MIT**

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
[**gRPC**](https://grpc.io/) is a modern open source high performance Remote Procedure Call (RPC) framework that can run in any environment. It can efficiently connect services in and across data centers with pluggable support for load balancing, tracing, health checking and authentication. It is also applicable in last mile of distributed computing to connect devices, mobile applications and browsers to backend services. While, gRPC as a protocol supports communication between services with JSON but it was designed to work well with the Protocol Buffers encoding format.

[**Protocol Buffers**](https://developers.google.com/protocol-buffers) are Google's language-neutral, platform-neutral, extensible mechanism for serializing structured data – think XML, but smaller, faster, and simpler. You define how you want your data to be structured once, then you can use special generated source code to easily write and read your structured data to and from a variety of data streams and using a variety of languages.

The protocol buffers for the FyrMesh are defined in a file ``fyrmesh/protos/fyrmesh.proto``.

**How does it work?**  

In gRPC, a client application can directly call a method on a server application on a different machine as if it were a local object, making it easier for you to create distributed applications and services. As in many RPC systems, gRPC is based around the idea of defining a service, specifying the methods that can be called remotely with their parameters and return types. On the server side, the server implements this interface and runs a gRPC server to handle client calls. On the client side, the client has a stub (referred to as just a client in some languages) that provides the same methods as the server.

gRPC clients and servers can run and talk to each other in a variety of environments - from servers inside Google to your own desktop - and can be written in any of gRPC’s supported languages. So, for example, you can easily create a gRPC server in Java with clients in Go, Python, or Ruby. In addition, the latest Google APIs will have gRPC versions of their interfaces, letting you easily build Google functionality into your applications.

gRPC Source files:
- [gRPC GitHub Repository](https://github.com/grpc/grpc)
- [gRPC-Go GitHub Repository](https://github.com/grpc/grpc-go)

For more information:
- [gRPC Introduction](https://grpc.io/docs/what-is-grpc/introduction/)
- [gRPC Core Concepts](https://grpc.io/docs/what-is-grpc/core-concepts/)
- [Protocol Buffers Overview](https://developers.google.com/protocol-buffers/docs/overview)
- [Working with Protocol Buffers](https://grpc.io/docs/what-is-grpc/introduction/#working-with-protocol-buffers)

### **Installation**
The FyrMesh library can be installed onto a Linux or a Windows system. Linux support is only intended for the Raspbian OS that is designed for the Raspberry Pi. 

#### 1. Prerequisites
Install Go 1.16 and Python 3.8
   - [Install Go on the Raspberry Pi](https://www.jeremymorgan.com/tutorials/raspberry-pi/install-go-raspberry-pi/)
   - [Install Go on Windows](https://golang.org/doc/install)
   - [Install Python on the Raspberry Pi](https://projects.raspberrypi.org/en/projects/generic-python-install-python3#:~:text=Open%20your%20web%20browser%20and,a%20download%20will%20start%20automatically.)
   - [Install Python on Windows](https://realpython.com/installing-python/)


Download this repository to some location on your system with one of the following methods. This location will be used to configure ``FYRMESHSRC`` in the next steps.
   - Download the latest release from the [releases](https://github.com/fyrwatch/fyrmesh/releases) (v0.2 or above).
   - Use the Git bash to clone the repository.  
  ``` 
  git clone https://github.com/fyrwatch/fyrmesh.git 
  ```
   - Use ``go get`` to download the package.
  ```
  go get -u github.com/fyrwatch/fyrmesh
  ```

At this stage, the project is only an application and not a fully managed service, so the set up of the cloud interface is upto the user. This involves having a **Google Cloud Platform** Project with billing enabled and a Firestore database created. An IAM Service Account with the role *Datastore User* must be created and JSON key for it must be generated. This key must exist in the ``FYRMESHCONFIG`` directory with the name *cloudconfig.json*

#### 2. Install FyrMesh
- Navigate into the ``/fyrmesh`` directory of the repository after downloading it.
- Open a terminal window in this directory and run the following command
```
python install.py
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
  export FYRMESHCONFIG=<path to the config file directory>
  export FYRMESHSRC=<path to the repository source files directory>
  export PATH=$PATH:<path to Go bin install directory>
  ``` 
  The recommended value for ``FYRMESHCONFIG`` is ``/home/pi/.fyrmesh``  
  The recommended value for go bin path extension is ``/home/pi/go/bin``
- Reboot the system to apply changes.

##### Windows
- Press the Start Menu and type *'env*'. 
- An option *'Edit the environment variables for your account'* will appear. Select it.
- A list of environment variables will be visible. Select the *'New'* option.
- Fill in the value for ``FYRMESHCONFIG``. The recommended value is ``C:\Users\<user>\.fyrmesh``.
- Repeat and fill in the value for ``FYRMESHSRC`` with the path to the install directory.
- PATH extension wouldn't be required because the Go installer for windows handles it.
- Reboot the systems to apply changes.

The installation of FyrMesh is now complete.  
Run the ``fyrcli`` command on the console to try it.

### Using the FyrCLI
Run the command ``fyrcli help`` to view the usage of the FyrCLI Application.
The output will be something along the lines of the following. 

```
A CLI Application to interact with the FyrMesh API. Powered by Cobra CLI, Golang and gRPC.

Usage:
  fyrcli [command]

Available Commands:
  boot        Boots a FyrMesh gRPC server.
  command     Sends a control command to the mesh.
  config      View configuration values of the FyrCLI.
  connect     Set the connection state of the control node.
  help        Help about any command
  nodelist    Displays the list of nodes connected to the mesh.
  observe     Observes the logstream from the ORCH server.
  ping        Pings the mesh or a node.
  scheduler   Sets the state of the Scheduler
  simulate    Starts a simulation of a Fire Event
  status      Displays the current status of the mesh.

Flags:
      --config string   config file (default is $HOME/.cli.yaml)
  -h, --help            help for fyrcli

Use "fyrcli [command] --help" for more information about a command.
```

You can also run ``fyrcli help <command>`` or ``fyrcli <command> --help`` to recieve the
usage instruction for any particular command.

### Starting the FyrMesh
Starting the FyrMesh involves booting ``LINK`` and ``ORCH`` services. After which commands can be called and they will perform as expected. The FyrMesh service can however only be run a device that is configured as mesh-controller, which is usually a ARM based controller like the Raspberry Pi.

Running the FyrCLI on Windows machines configured as a mesh-observer is simply useful to send commands remotely and observing the behaviour of the mesh-controller.

**How to boot the services?**

To boot the ``ORCH`` and ``LINK`` services, run the following commands.
```
fyrcli boot --server LINK
fyrcli boot --server ORCH
```

The ``LINK`` server must **always** be booted before the ``ORCH`` server to avoid gRPC errors.
