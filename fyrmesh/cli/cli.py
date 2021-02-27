import os
import sys
import json
import time
import click
import http.client

FYRMESH_SERVER = http.client.HTTPConnection('localhost', 8888)

def isConnected() -> bool:
    try:
        FYRMESH_SERVER.request('GET', '/', '{}')
        response = json.loads(FYRMESH_SERVER.getresponse().read())
        del response
        return True
    except Exception:
        return False

def checkConnection():
    if not isConnected():
        click.echo("The FyrMesh server is not running. Run 'fyrmesh boot' to start the server.")
        sys.exit()


@click.group()
def cli():
    pass

@cli.command()
def boot():
    if isConnected():
        click.echo("The FyrMesh server is already running.")
        sys.exit()

    directory = os.path.dirname(os.path.realpath(__file__))
    serverscript = os.path.join(directory, '..', 'server','server.py')
    os.system(f"start python {serverscript}")

@cli.command()
def activate():
    checkConnection()

    FYRMESH_SERVER.request('GET', '/activate', '{}')
    response = json.loads(FYRMESH_SERVER.getresponse().read())

    click.echo(f"{response['message']}")

@cli.command()
def deactivate():
    checkConnection()

    FYRMESH_SERVER.request('GET', '/deactivate', '{}')
    response = json.loads(FYRMESH_SERVER.getresponse().read())

    click.echo(f"{response['message']}")

@cli.command()
def status():
    checkConnection()

    FYRMESH_SERVER.request('GET', '/status', '{}')
    response = json.loads(FYRMESH_SERVER.getresponse().read())

    click.echo(f"{response['message']}")
