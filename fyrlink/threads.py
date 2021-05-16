"""
===========================================================================
Copyright (C) 2020 Manish Meganathan, Mariyam A.Ghani. All Rights Reserved.
 
This file is part of the FyrMesh library.
No part of the FyrMesh library can not be copied and/or distributed 
without the express permission of Manish Meganathan and Mariyam A.Ghani
===========================================================================
FyrMesh FyrLINK module

A module that contains implementations for a class KillableThread that 
allows a thread to be killed with a command, a functionality that is 
not availabe in the regular threading module.
===========================================================================
"""

import threading
import ctypes

class KillableThread(threading.Thread):
    """ Custom thread class that inherits from threading.Thread and supports 
    the ability to kill the created thread after it has been started """

    def __init__(self, name: str, *args, **kwargs):
        threading.Thread.__init__(self, *args, **kwargs)
        self.name = name

    def get_id(self):
        if hasattr(self, '_thread_id'):
            return self._thread_id
        for id, thread in threading._active.items():
            if thread is self: return id

    def kill(self):
        thread_id = self.get_id() 
        res = ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id, ctypes.py_object(SystemExit)) 

        if res > 1:
            ctypes.pythonapi.PyThreadState_SetAsyncExc(thread_id, 0) 
            print('Thread Kill Failed') 
