"""
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.
 
This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed 
without the express permission of Manish Meganathan and Mariyam A.Ghani
===========================================================================
FyrMesh FyrLINK module

A module that contains implementations for thread workers and 
the tools that run within them. The serial port object and the 
read/write queues are defined here as well a class KillableThread 
that allows a thread to be killed with a command, a functionality 
that is not availabe in the regular threading module.
===========================================================================
"""

import json
import queue
import serial
import ctypes
import threading
from fyrlink.parsers import logtime

# Define the log and command queues
logqueue = queue.Queue()
commandqueue = queue.Queue()

# Define a sync lock for the logger
loglock = threading.Lock()

# Define the Serial Port interface
serialport = serial.Serial(
    port='/dev/ttyAMA0',
    baudrate='115200',
    parity=serial.PARITY_NONE,
    stopbits=serial.STOPBITS_ONE,
    bytesize=serial.EIGHTBITS,
    timeout=None
)

class KillableThread(threading.Thread):
    """ Custom thread class that inherits from threading.Thread and supports 
    the ability to kill the created thread after it has been started """

    def __init__(self, name: str, *args, **kwargs):
        threading.Thread.__init__(self, *args, **kwargs)
        self.name = name

    def get_id(self):
        if hasattr(self, '_thread_id'):
            return self._thread_id
        for id, thread in threading._active.items():
            if thread is self: return id

    def kill(self):
        thread_id = self.get_id() 
        res = ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id, ctypes.py_object(SystemExit)) 

        if res > 1:
            ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id, 0)
            logqueue.put({
                "source": "LINK", 
                "type": "serverlog", 
                "time": logtime(), 
                "log": f"Attempt to kill the {self.name} thread failed!",
                "metadata": {}
            }) 


def readfromqueue(somequeue: queue.Queue):
    """ A function that takes a Queue object and returns an element from it. 
    Returns a None object if there is no element in the Queue """
    if not somequeue.empty():
        try:
            member = somequeue.get(block=False)
            return member
        except queue.Empty:
            return None
    else:
        return None

def writer():
    """ A threadworker function that reads from the commandqueue and writes to the serial port """
    while True:
        command = readfromqueue(commandqueue)
        serialport.write(json.dumps(command).encode('ascii')) if command else None

def reader():
    """ A threadworker function that reads from the serial port, parses it and adds it to the logqueue """
    from fyrlink.parsers import parse

    while True:
        rxdata = serialport.read_until()
        parsed = parse(rxdata)
        logqueue.put(parsed) if parsed else None

def logger():
    """ A threadworker function that reads from the logqueue and prints it to the console """
    while True:
        with loglock:
            log = readfromqueue(logqueue)
            print(log) if log else None
