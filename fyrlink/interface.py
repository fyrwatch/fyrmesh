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
import grpc
import time
import threading
import concurrent.futures as futures
import proto.fyrmesh_pb2 as fyrmesh_pb2
import proto.fyrmesh_pb2_grpc as fyrmesh_pb2_grpc

from fyrlink.parsers import logtime
from fyrlink.workers import writequeue, readqueue, loglock
from fyrlink.workers import reader, writer, logger, readfromqueue
from fyrlink.threads import KillableThread

# Create lock synchronisation object
threadlock = threading.Lock()

class Interface(fyrmesh_pb2_grpc.InterfaceServicer):
    """ Class that implements the Interface gRPC Server """

    def Read(self, request, context):
        """ A method that implements the 'Read' runtime of the Interface 
        server. Collects log messages from the read queue and streams them 
        to the gRPC Interface client. """

        if request.message == "start-stream-read":
            with loglock:       
                while True:
                    message = readfromqueue(readqueue)

                    if message:
                        yield fyrmesh_pb2.InterfaceLog(
                            logsource=message['source'], 
                            logtime=message['time'], 
                            logmessage=message['log'])
                    else:
                        pass
        else:
            while True:
                time.sleep(2)
                yield fyrmesh_pb2.InterfaceLog(
                    logsource="LINK", 
                    logtime=logtime(), 
                    logmessage="invalid read stream initiation code")

        
    def Write(self, request, context):
        """ A method that implements the 'Write' runtime of the Interface 
        server. Puts the command recieved into the write queue with the
        appropriate structure. """

        command = request.command
        writequeue.put({"type": "controlcommand", "command": command})
        return fyrmesh_pb2.Acknowledge(success=True, error="nil")


def grpc_serve():
    """ A function that sets up and serves the gRPC LINK Interface Server """
    # Create a gRPC server object
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    # Register the server as an Interface server
    fyrmesh_pb2_grpc.add_InterfaceServicer_to_server(Interface(), server)

    # TODO: Extract port information from the config file.
    # Setup the server listening port and start it.
    server.add_insecure_port(f'[::]:50000')
    server.start()

    # Log the start of the server.
    readqueue.put({"source": "LINK", "time": logtime(), "log": "Interface Link gRPC Server started on http://localhost:50000"})

    # Server will wait indefinitely for termination
    server.wait_for_termination()


if __name__ == "__main__":
    # Define the IO thread workers that run concurrently
    readerthread = KillableThread(name="reader", target=reader, daemon=True)
    writerthread = KillableThread(name="writer", target=writer, daemon=True)
    loggerthread = KillableThread(name="logger", target=logger, daemon=True)

    # Start the IO thread workers
    readerthread.start()
    writerthread.start()
    loggerthread.start()

    try:
        # Start the gRPC server
        grpc_serve()
    
    except KeyboardInterrupt:
        # Kill the IO thread workers
        readerthread.kill()
        writerthread.kill()
        loggerthread.kill()
    
    # Exit without handling regular exit runtimes such as printing tracebacks
    os._exit(1)
