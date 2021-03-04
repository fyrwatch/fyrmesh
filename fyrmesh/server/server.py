import rpyc
import threading

meshuptime = 0
meshstatus = False
meshlock = threading.Lock()

def countuptime():
    import time
    global meshstatus, meshuptime

    while 1:
        with meshlock:
            meshuptime = meshuptime+1 if meshstatus else meshuptime
        time.sleep(1)

class MeshService(rpyc.Service):
    def on_connect(self, conn):
        pass

    def on_disconnect(self, conn):
        pass

    @property
    def exposed_activate(self):
        global meshstatus

        with meshlock:
            meshstatus = True

        print("mesh activated")
        return "mesh activated"

    @property
    def exposed_deactivate(self):
        global meshstatus

        with meshlock:
            meshstatus = False

        print("mesh deactivated")
        return "mesh deactivated"

    @property
    def exposed_status(self):
        import json
        global meshstatus, meshuptime

        with meshlock:
            message = { 
                "meshstatus": "ACTIVE" if meshstatus else "INACTIVE",
                "meshuptime": f"{meshuptime}s"
            }

        return json.dumps(message)


if __name__ == "__main__":
    counter = threading.Thread(target=countuptime, daemon=True)
    counter.start()

    from rpyc.utils.server import ThreadedServer
    server = ThreadedServer(MeshService, port=18000)
    server.start()