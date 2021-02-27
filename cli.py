import os
import sys
import json
import click
import http.client

fyrmesh_server = http.client.HTTPConnection('localhost', 8888)

def isConnected() -> bool:
    try:
        fyrmesh_server.request('GET', '/', '{}')
        response = json.loads(fyrmesh_server.getresponse().read())
        connection = True
    except Exception:
        connection = False

    return connection

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
    serverscript = os.path.join(directory, 'server.py')
    os.system(f"start python {serverscript}")

    message = "The FyrMesh server booted succefully" if isConnected() else "The FyrMesh server failed to boot"
    click.echo(message)

@cli.command()
def activate():
    checkConnection()

    fyrmesh_server.request('GET', '/activate', '{}')
    response = json.loads(fyrmesh_server.getresponse().read())

    click.echo(f"{response['message']}")

@cli.command()
def deactivate():
    checkConnection()

    fyrmesh_server.request('GET', '/deactivate', '{}')
    response = json.loads(fyrmesh_server.getresponse().read())

    click.echo(f"{response['message']}")

@cli.command()
def status():
    checkConnection()

    fyrmesh_server.request('GET', '/status', '{}')
    response = json.loads(fyrmesh_server.getresponse().read())

    click.echo(f"{response['message']}")
