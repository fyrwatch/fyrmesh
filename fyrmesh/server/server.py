import json
import time
import threading
from bottle import run, route

synclock = threading.Lock()
meshstatus = False
meshuptime = 0

def countuptime():
    global meshstatus, meshuptime

    while 1:
        with synclock:
            meshuptime = meshuptime+1 if meshstatus else meshuptime
        time.sleep(1)

@route('/')
def home():
    return json.dumps({
        "message": "mesh homepage"
    })

@route('/activate')
def activatemesh():
    global meshstatus

    with synclock:
        meshstatus = True

    return json.dumps({
        "message": "mesh activated"
    })

@route('/deactivate')
def deactivatemesh():
    global meshstatus

    with synclock:
        meshstatus = False
    
    return json.dumps({
        "message": "mesh deactivated"
    })


@route('/status')
def meshstatusfunc():
    global meshstatus, meshuptime

    with synclock:
        status = "ACTIVE" if meshstatus else "INACTIVE"
        count = f"{meshuptime}s"

    return json.dumps({
        "message": {
            "meshstatus": status,
            "meshuptime": count
        }
    })

if __name__ == "__main__":
    counter = threading.Thread(target=countuptime, daemon=True)
    counter.start()
    run(port=8888)