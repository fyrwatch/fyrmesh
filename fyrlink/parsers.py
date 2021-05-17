"""
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.
 
This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed 
without the express permission of Manish Meganathan and Mariyam A.Ghani
===========================================================================
FyrMesh FyrLINK module

A module that contains functions that are used to parse and messages
recieved from the serial port and restructure them appropriately.
===========================================================================
"""

import json
import datetime

logtime = lambda: datetime.datetime.utcnow().isoformat()

def dictparse(data):
    """ A function that attempts to parse the receieved data as a JSON dictonary"""
    parseddata = json.loads(data)

    if parseddata['type'] == "meshlog":
        return meshlogparse(parseddata)  
    else: 
        return {
            "source": "MESH", 
            "type": "message", 
            "time": logtime(), 
            "log": f"{parseddata}", 
            "metadata": {
                "format": "dict", 
                "type": parseddata['type']
            }
        }

def strparse(data):
    """ A function that attempts to parse the receieved data as an ASCII string """
    try:
        parsed = data.decode('ascii')
        return {
            "source": "MESH", 
            "type": "message", 
            "time": logtime(), 
            "log": f"{parsed.rstrip()}", 
            "metadata": {
                "format": "str", 
                "type": "unknown"
            }
        } if parsed else None

    except:
        return None

def meshlogparse(meshlog: dict):
    """ A function that parses a meshlog dictionary and adds it to the logqueue with the approptiate formatting """

    meshlogdata = meshlog['logdata']
    meshlogtype = meshlogdata['type']
    logmessage = {"source": "MESH", "time": logtime()}

    if meshlogtype == "newconnection":
        logmessage.update({
            "type": "newconnection", 
            "log": "new node on mesh", 
            "metadata": {
                "node": meshlogdata['newnode']
            }
        })
    
    elif meshlogtype == "changedconnection":
        logmessage.update({
            "type": "changedconnection", 
            "log": "mesh connections have changed", 
            "metadata": {}
        })

    elif meshlogtype == "nodetimeadjust":
        logmessage.update({
            "type": "nodetimeadjust", 
            "log": "node time was adjusted to sync with the mesh", 
            "metadata": {
                "offset": meshlogdata['offset']
            }
        })

    elif meshlogtype == "handshake-rxack":
        logmessage.update({
            "type": "handshake", 
            "log": "handshake completed with a node", 
            "metadata": {
                "node": meshlogdata['node']
            }
        })

    elif meshlogtype == "sensordata":
        logmessage.update({
            "type": "sensordata", 
            "log": "sensor data received", 
            "metadata": {
                "ping": meshlogdata['ping'],
                "node": meshlogdata['node'],
                "sensors": meshlogdata['sensors']
            }
        })

    elif meshlogtype == "configdata":
        logmessage.update({
            "type": "configdata", 
            "log": "config data received", 
            "metadata": {
                "ping": meshlogdata['ping'],
                "node": meshlogdata['node'],
                "config": meshlogdata['config']
            }
        })

    elif meshlogtype == "controlconfigdata":
        logmessage.update({
            "type": "controlconfig",
            "log": "control node config data received",
            "metadata": {
                "nodeID": meshlogdata["config"]["NODEID"],
                "config": meshlogdata["config"]
            }
        })

    elif meshlogtype == "controlnodelist":
        logmessage.update({
            "type": "nodelist",
            "log": "mesh node list received",
            "metadata": {
                "nodelist": meshlogdata["nodelist"]
            }
        })

    elif meshlogtype == "messagerx":
        logmessage.update({
            "type": "message", 
            "log": f"{meshlogdata['message']}", 
            "metadata": {
                "format": "str", 
                "type": meshlogdata['rxtype']
                }
            })

    else:
        logmessage.update({
            "type": "message", 
            "log": f"{meshlogdata}",
            "metadata": {
                "format": "dict",
                "type": meshlogtype
                }
            })

    return logmessage

def parse(data: bytes):
    """ A function that parses a byte string data into an appropriate mesh message """
    try: 
        return dictparse(data)
    except:
        return strparse(data)
