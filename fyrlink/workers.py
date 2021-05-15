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
read/write queues are defined here as well.
===========================================================================
"""

import json
import queue
import serial
import threading

# Define the read and write queues
readqueue = queue.Queue()
writequeue = queue.Queue()

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
    """ A threadworker function that reads from the writequeue and writes to the serial port """
    while True:
        command = readfromqueue(writequeue)
        serialport.write(json.dumps(command).encode('ascii')) if command else None

def reader():
    """ A threadworker function that reads from the serial port, parses it and adds it to the readqueue """
    from fyrlink.parsers import parse

    while True:
        rxdata = serialport.read_until()
        parsed = parse(rxdata)
        readqueue.put(parsed) if parsed else None

def logger():
    """ A threadworker function that reads from the readqueue and prints it to the console """
    while True:
        with loglock:
            log = readfromqueue(readqueue)
            print(log) if log else None
