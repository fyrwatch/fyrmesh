"""
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.
 
This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed 
without the express permission of Manish Meganathan and Mariyam A.Ghani
===========================================================================
FyrMesh FyrLINK server
===========================================================================
"""

import os
import sys
import json
import grpc
import time
import threading
import concurrent.futures as futures
import proto.fyrmesh_pb2 as fyrmesh_pb2
import proto.fyrmesh_pb2_grpc as fyrmesh_pb2_grpc

from fyrlink.parsers import logtime
from fyrlink.workers import commandqueue, logqueue, loglock
from fyrlink.workers import reader, writer, logger, readfromqueue
from fyrlink.workers import KillableThread

# Create lock synchronisation object
threadlock = threading.Lock()

class Interface(fyrmesh_pb2_grpc.InterfaceServicer):
    """ Class that implements the Interface gRPC Server """

    def Read(self, request, context):
        """ A method that implements the 'Read' runtime of the Interface 
        server. Collects log messages from the read queue and streams them 
        to the gRPC Interface client. """

        if request.triggermessage == "start-stream-read":
            with loglock:       
                while True:
                    message = readfromqueue(logqueue)

                    if message:
                        yield fyrmesh_pb2.ComplexLog(
                            logsource=message['source'], 
                            logtype=message['type'],
                            logtime=message['time'], 
                            logmessage=message['log'], 
                            logmetadata=message['metadata']
                        )

                    else:
                        pass
        else:
            while True:
                time.sleep(2)
                yield fyrmesh_pb2.Complexlog(
                    logsource="LINK", 
                    logtype="protolog",
                    logtime=logtime(), 
                    logmessage="(error) invalid read stream initiation code", 
                    logmetadata={
                        "server": "LINK",
                        "service": "Read",
                        "error": "nil"
                })

        
    def Write(self, request, context):
        """ A method that implements the 'Write' runtime of the Interface 
        server. Puts the command recieved into the write queue with the
        appropriate structure. """

        command = request.command
        metadata = request.metadata

        try:
            commandqueue.put({"type": "controlcommand", "command": command, **metadata})
            logqueue.put({
                "source": "LINK", "type": "protolog", "time": logtime(), 
                "log": f"(success) command '{command}' written to control node successfully",
                "metadata": {
                    "server": "LINK", 
                    "service": "Write",
                    "error": "nil"
            }})

        except Exception as e:
            logqueue.put({
                "source": "LINK", "type": "protolog", "time": logtime(), 
                "log": f"(failure) command '{command}' failed to be written to control node.",
                "metadata": {
                    "server": "LINK",
                    "service": "Write",
                    "error": str(e)
            }})
            return fyrmesh_pb2.Acknowledge(success=False, error=str(e))

        return fyrmesh_pb2.Acknowledge(success=True, error="nil")


def grpc_serve():
    """ A function that sets up and serves the gRPC LINK Interface Server """
    # Create a gRPC server object
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    # Register the server as an Interface server
    fyrmesh_pb2_grpc.add_InterfaceServicer_to_server(Interface(), server)

    # Retrieve the FYRMESHCONFIG env var
    configpath = os.environ.get("FYRMESHCONFIG")
    if not configpath:
        logqueue.put({
            "source": "LINK", 
            "type": "serverlog", 
            "time": logtime(), 
            "log": "(error) could not read config. 'FYRMESHCONFIG' env variable is not set",
            "metadata": {}
        })
        sys.exit()

    # Read the config file
    configfilepath = os.path.join(configpath, "config.json")
    with open(configfilepath) as configfile:
        configdata = json.load(configfile)

    # Setup the server listening port and start it.
    port = configdata['services']['LINK']['port']
    server.add_insecure_port(f'[::]:{port}')
    server.start()

    # Log the start of the server.
    logqueue.put({
        "source": "LINK", 
        "type": "serverlog", 
        "time": logtime(), 
        "log": "(startup) interface link grpc server started on http://localhost:50000",
        "metadata": {}
    })

    # Server will wait indefinitely for termination
    server.wait_for_termination()


if __name__ == "__main__":
    # Define the IO thread workers that run concurrently
    readerthread = KillableThread(name="reader", target=reader, daemon=True)
    writerthread = KillableThread(name="writer", target=writer, daemon=True)
    #loggerthread = KillableThread(name="logger", target=logger, daemon=True)

    # Start the IO thread workers
    readerthread.start()
    writerthread.start()
    #loggerthread.start()

    try:
        # Start the gRPC server
        grpc_serve()
    
    except KeyboardInterrupt:
        # Kill the IO thread workers
        readerthread.kill()
        writerthread.kill()
        #loggerthread.kill()
    
    # Exit without handling regular exit runtimes such as printing tracebacks
    os._exit(1)
