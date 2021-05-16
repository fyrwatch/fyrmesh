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
    parsed = json.loads(data)

    if parsed['type'] == "meshlog":
        return meshlogparse(parsed)  
    else: 
        return {"source": "MESH", "time": logtime(), "log": f"(dict) {parsed}"}

def strparse(data):
    """ A function that attempts to parse the receieved data as an ASCII string """
    try:
        parsed = data.decode('ascii')
        return {"source": "MESH", "time": logtime(), "log": f"(str) {parsed.rstrip()}"} if parsed else None
    except:
        return None

def meshlogparse(meshlog: dict):
    """ A function that parses a meshlog dictionary and adds it to the logqueue with the approptiate formatting """

    meshlogdata = meshlog['logdata']
    meshlogtype = meshlogdata['type']
    logmessage = {"source": "MESH", "time": logtime()}

    if meshlogtype == "newconnection":
        logmessage.update({"log": f"(newconnection) node={meshlogdata['newnode']}"})
    
    elif meshlogtype == "changedconnection":
        logmessage.update({"log": f"(changedconnection)"})

    elif meshlogtype == "nodetimeadjust":
        logmessage.update({"log": f"(nodetimeadjust) offset={meshlogdata['offset']}"})

    elif meshlogtype == "sensordata":
        logmessage.update({"log": f"(sensordata) poll={meshlogdata['poll']} node={meshlogdata['node']} sensors={meshlogdata['sensors']}"})

    elif meshlogtype == "handshake-rxack":
        logmessage.update({"log": f"(handshake-rxack) node={meshlogdata['node']}"})

    elif meshlogtype == "messagerx":
        logmessage.update({"log": f"(messagerx) type={meshlogdata['rxtype']} message={meshlogdata['message']}"})

    else:
        logmessage.update({"log": f"(unknowntype) type={meshlogtype})"})

    return logmessage

def parse(data: bytes):
    """ A function that parses a byte string data into an appropriate mesh message """
    try: 
        return dictparse(data)
    except:
        return strparse(data)
