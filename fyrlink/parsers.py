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
FyrMesh FyrLINK module

A module that contains functions that are used to parse and messages
recieved from the serial port and restructure them appropriately.
===========================================================================
"""

import json
import datetime

logtime = lambda: datetime.datetime.utcnow().isoformat().split(".")[0]

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

def deepserialize(data: dict):
    """ A function that serializes a dictionary into a string with a format akin to 'key1-value1=key2-value2..'"""
    strdata = ""
    for key, value in data.items():
        strdata = strdata + f"{key}-{value}="

    strdata = strdata.strip("=")
    return strdata

def meshlogparse(meshlog: dict):
    """ A function that parses a meshlog dictionary and adds it to the logqueue with the approptiate formatting """

    meshlogdata = meshlog['logdata']
    meshlogtype = meshlogdata['type']
    logmessage = {"source": "MESH", "time": logtime()}

    if meshlogtype == "meshsync":
        logmessage.update({
            "type": "meshsync",
            "log": "mesh synchronization event",
            "metadata": {
                "synctype": meshlogdata['sync']
            }
        })

    elif meshlogtype == "nodesync":
        logmessage.update({
            "type": "nodesync", 
            "log": "node time synchronized", 
            "metadata": {
                "offset": str(meshlogdata['offset'])
            }
        })

    elif meshlogtype == "handshake-rxack":
        logmessage.update({
            "type": "handshake", 
            "log": "node handshaked", 
            "metadata": {
                "node": str(meshlogdata['node'])
            }
        })

    elif meshlogtype == "sensordata":
        logmessage.update({
            "type": "sensordata", 
            "log": "sensordata acquired", 
            "metadata": {
                "ping": meshlogdata['ping'],
                "node": str(meshlogdata['node']),
                "sensors": deepserialize(meshlogdata['sensors'])
            }
        })

    elif meshlogtype == "configdata":
        logmessage.update({
            "type": "configdata", 
            "log": "configdata acquired", 
            "metadata": {
                "ping": meshlogdata['ping'],
                "node": str(meshlogdata['node']),
                "config": deepserialize(meshlogdata['config'])
            }
        })

    elif meshlogtype == "controlconfigdata":
        logmessage.update({
            "type": "ctrldata",
            "log": "controlnode config acquired",
            "metadata": {
                "nodeID": str(meshlogdata["config"]["NODEID"]),
                "config": deepserialize(meshlogdata['config'])
            }
        })

    elif meshlogtype == "controlnodelist":
        logmessage.update({
            "type": "nodelist",
            "log": "mesh nodelist acquired",
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
