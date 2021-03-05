import os
import sys
import rpyc
import click

mesh = None
MESH_CONNECTED = False

try:
    conn = rpyc.connect("localhost", 18000)
    mesh = conn.root
    MESH_CONNECTED = True
except Exception:
    MESH_CONNECTED = False

def checkConnection():
    if not MESH_CONNECTED:
        click.echo("The FyrMesh server is not running. Run 'fyrmesh boot' to start the server.")
        sys.exit()


@click.group()
def cli():
    pass

@cli.command()
def boot():
    if MESH_CONNECTED:
        click.echo("The FyrMesh server is already running.")
        sys.exit()

    directory = os.path.dirname(os.path.realpath(__file__))
    serverscript = os.path.join(directory, '..', 'server','server.py')
    os.system(f"start python {serverscript}")

@cli.command()
def activate():
    checkConnection()

    response = mesh.activate
    click.echo(response)

@cli.command()
def deactivate():
    checkConnection()

    response = mesh.deactivate
    click.echo(response)

@cli.command()
def status():
    checkConnection()

    response = mesh.status
    click.echo(response)
