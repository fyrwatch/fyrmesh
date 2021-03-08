import os
import sys
import rpyc
import json
import click

mesh = None
MESH_CONNECTED = False
serverhost = None
serverport = None

def connectserver(host: str, port: int) -> bool:
    """A function that attempts to connect to the server and returns 
    whether the attempt was successful"""

    global mesh

    try:
        conn = rpyc.connect(host, port)
        mesh = conn.root
        return True
    except Exception:
        return False

def loadserver():
    """A function that loads the server values from the configuration file and connects to it.
    If the configuration file is not found, it is created.
    """
    
    global MESH_CONNECTED

    if not os.path.isfile(configfile):
        with open(configfile, "w") as fp:
            serverdata = {"host": "localhost", "port": 18000}
            json.dump(serverdata, fp)

    else:
        with open(configfile) as fp:
            serverdata = json.load(fp)
        
    serverhost = serverdata['host']
    serverport = serverdata['port']
    MESH_CONNECTED = connectserver(serverhost, serverport)

directory = os.path.dirname(os.path.realpath(__file__))
serverscript = os.path.join(directory, '..', 'server','server.py')
configfile = os.path.join(directory, 'cliconfig.json')

loadserver()

def checkconnection():
    """A function that checks if the connection is already established"""
    if not MESH_CONNECTED:
        bootmessage = ("""The FyrMesh server is not connected.Run 'fyrmesh boot' to start the server on this 
        machine or set the host and port for a remote machine with 'fyrmesh setconfig' command""")
        click.echo(bootmessage)
        sys.exit()

@click.group()
def cli():
    """FyrMesh CLI to interact with the mesh of sensor nodes"""
    pass

@cli.command()
def boot():
    """boots up the FyrMesh server"""
    if MESH_CONNECTED:
        click.echo("The FyrMesh server is already running.")
        sys.exit()
    
    os.system(f"lxterminal -e 'python3 {serverscript}' &")

@cli.command()
def connect():
    """connects the server to the mesh"""
    checkconnection()
    response = mesh.connectmesh()
    click.echo(response)

@cli.command()
def disconnect():
    """disconnects the server from the mesh"""
    checkconnection()
    response = mesh.disconnectmesh()
    click.echo(response)

@cli.command()
def status():
    """returns the current status of the mesh"""
    checkconnection()
    response = mesh.status()
    click.echo(response)

@cli.command()
def poll():
    """polls the sensors on all nodes on the mesh"""
    checkconnection()
    mesh.pollmeshsensors()

@cli.command()
def reloadconfig():
    """reloads the FyrMesh connection by realoading the configuration values"""
    loadserver()

@cli.command()
@click.option('--host', default="localhost", help='IP address of the host')
@click.option('--port', default=18000, help='port of the host')
def setconfig(host, port):
    """set the FyrMesh server host and port values to a configuration file"""
    with open(configfile, "w") as fp:
        serverdata = {"host": host, "port": port}
        json.dump(serverdata, fp)
