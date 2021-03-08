# FyrMesh
![FyrMesh Banner](fyrmesh.png)
## A Python package for mesh orchestration and control with a built-in CLI and RPC-based mesh handling server.

**Version: 0.1.0**  
**Platform: Raspberry Pi (Linux)**  
**Language: Python 3.7**

### **Contributors**
- **Manish Meganathan**
- **Mariyam A. Ghani**

This repository is a python package that includes: 
- **An RPC Server Program**  
  An RPC Service that runs the orchestartion interface between the mesh and the Raspberry Pi while exposing the methods of the service as an RPC object that is hosted on Raspberry Pi and can be connected through the localnet.
- **A Command Line Interface library**  
  An interface library that exposes commands as a CLI built using the Click framework. It connects to the RPC server and can be spun on up on machine on the local network and configured.
- **A Firbase Cloud Interface library**  
  An interface library that communicates with the Firebase backend and other cloud services.
- **A Mesh Orchestration library**  
  A library of custom thread runtimes, data classes, mesh messaging protocol runtimes, etc. that play a role in the orchestation of the mesh.


### **Installation** 
The ``fyrmesh`` package and CLI can be installed by first ``cd`` into the repository and running the command:
```
pip3 install -e .
```
The ``-e`` flag sets the installation to *editable*. This is reserved just for development use. Similarly the use between ``pip3`` and ``pip`` is dependant on the configuration of the machine.

The fyrmesh CLI is now installed. Run ``fyrmesh`` on the terminal to see the help document. If an 'entry point not found' error is thrown, do the following:  
- Run the following to open the ``.bash_profile`` file.
  ```
  nano ~/.bash_profile
  ``` 
- Add the following line to the file.
  ```
  export PATH=~/.local/bin:$PATH
  ```
- Run the following on a terminal
  ```
  source ~/.bash_profile
  ```
The last command might have to be run for each new terminal window depending on the machine configuration.