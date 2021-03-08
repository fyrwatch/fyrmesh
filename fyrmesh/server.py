import rpyc
import json
import queue
import serial
import datetime
import threading

meshconnected = True
readerthread = None
printerthread = None
commandthread = None

threadlock = threading.Lock()
printqueue = queue.Queue()
commandqueue = queue.Queue()

serialport = serial.Serial(
    port='/dev/ttyAMA0',
    baudrate='115200',
    parity=serial.PARITY_NONE,
    stopbits=serial.STOPBITS_ONE,
    bytesize=serial.EIGHTBITS,
    timeout=None)

def logtime():
    """A function that returns the current time in an ISO string"""
    return datetime.datetime.utcnow().isoformat()

def parse_meshlog(meshlog: dict):
    """A function that parses the meshlog dictionary, creates a server log string and
    adds the serverlog into the printqueue"""

    logdata = meshlog['logdata']
    logtype = logdata['type']

    if logtype == "newconnection":
        logmessage = f"[MESHLOG][{logtime()}][newconnection] {logdata['newnode']}"
    
    elif logtype == "changedconnection":
        logmessage = f"[MESHLOG][{logtime()}][changedconnection]"

    elif logtype == "nodetimeadjust":
        logmessage = f"[MESHLOG][{logtime()}][nodetimeadjust] {logdata['offset']}]"

    elif logtype == "sensordata":
        logmessage = f"[MESHLOG][{logtime()}][sensordata] poll-{logdata['poll']} node-{logdata['node']} {logdata['sensors']}"

    elif logtype == "handshake-rxack":
        logmessage = f"[MESHLOG][{logtime()}][handshake-rxack] {logdata['node']}"

    elif logtype == "messagerx":
        logmessage = f"[MESHLOG][{logtime()}][messagerx] type-{logdata['rxtype']} {logdata['message']}"

    else:
        logmessage = f"[MESHLOG][{logtime()}][unknowntype] {logtype}]"

    with threadlock:
        printqueue.put(logmessage)

def dictparse(rxdata):
    """A function that attempts to parse the receieved data as a dictonary"""
    finaldata = json.loads(rxdata)
    parse_meshlog(finaldata) if finaldata['type'] == "meshlog" else print(f"[MESHDATA][{logtime()}][dict] {finaldata}")

def strparse(rxdata):
    """A function that attempts to parse the receieved data as string"""
    try:
        finaldata = rxdata.decode('ascii')
        print(f"[MESHDATA][{logtime()}][str] {finaldata.rstrip()}") if finaldata else None
    except:
        del rxdata

def readworker():
    """A threadworker function that keeps reading the serial port and attempts to parse it """
    while True:
        rxdata = serialport.read_until()

        try: 
            dictparse(rxdata)
        except:
            strparse(rxdata)

def commandworker():
    """A threadworker function that keeps checking the commandqueue for commands and writes
    them to the serial port"""
    while True:
        if not commandqueue.empty():
            try:
                command = commandqueue.get(block=False)
                serialport.write(json.dumps(command).encode('ascii'))
            except queue.Empty:
                pass
        else:
            pass

def printworker():
    """A threadworker function that keeps checking the printqueue for new messages to log
    prints them to stdout"""
    while True:
        if not printqueue.empty():
            try:
                log = printqueue.get(block=False)
                print(log)
            except queue.Empty:
                pass
        else:
            pass

class MeshService(rpyc.Service):
    """An RPyC Service Class to for the Mesh"""
    def on_connect(self, conn):
        """A function that runs when a client connects to this service"""
        pass

    def on_disconnect(self, conn):
        """A function that runs when a client disconnect from this service"""
        pass

    def exposed_connectmesh(self):
        """A function that sends the 'connection-on' command to controlnode"""
        global meshconnected
        with threadlock:
            meshconnected = True

        command = {"type": "controlcommand", "command": "connection-on"}
        commandqueue.put(command)
        printqueue.put(f"[SERVER][{logtime()}] MESH CONNECTED")
        return "mesh connected"

    def exposed_disconnectmesh(self):
        """A function that sends the 'connection-off' command to controlnode"""
        global meshconnected
        with threadlock:
            meshconnected = False

        command = {"type": "controlcommand", "command": "connection-off"}
        commandqueue.put(command)
        printqueue.put(f"[SERVER][{logtime()}] MESH DISCONNECTED")
        return "mesh disconnected"

    def exposed_pollmeshsensors(self):
        """A function that sends the 'readsensors-mesh' command to controlnode"""
        command = {"type": "controlcommand", "command": "readsensors-mesh"}
        commandqueue.put(command)
        printqueue.put(f"[SERVER][{logtime()}] MESH SENSORS POLLED")

    def exposed_status(self):
        """A function that returns whether the service is currently connected to the controlnode"""
        return f"meshconnected {meshconnected}"


if __name__ == "__main__":
    readerthread = threading.Thread(target=readworker, daemon=True)
    printerthread = threading.Thread(target=printworker, daemon=True)
    commandthread = threading.Thread(target=commandworker, daemon=True)

    readerthread.start()
    printerthread.start()
    commandthread.start()

    printqueue.put(f"[SERVER][{logtime()}] FyrMesh Server Started http://localhost/18000")
    server = rpyc.utils.server.ThreadedServer(MeshService, port=18000)
    server.start()