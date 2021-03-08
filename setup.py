#!/usr/bin/env python3
"""
*************************************************************
* Copyright (C) 2020 Manish Meganathan 757manish@gmail.com
* 
* This file is part of FyrWatch and FyrMesh.
* 
* FyrWatch can not be copied and/or distributed without the 
* express permission of Manish Meganathan and Mariyam A.Ghani
*************************************************************
"""

from setuptools import setup

setup(
    name="fyrmesh",
    version="0.1.0",
    description='Fyrmesh Orchestration Service and API',
    url='https://github.com/fyrwatch/fyrmesh',
    packages=['fyrmesh'],
    package_data = {},
    test_suite = '',
    install_requires=[
        'click',
        'pyserial',
        'rpyc'
    ],
    entry_points = {
        'console_scripts': ['fyrmesh=fyrmesh.cli:cli']
    },
    classifiers=[
        # Get strings from
        # http://pypi.python.org/pypi?%3Aaction=list_classifiers
        'Programming Language :: Python :: 3',
        'Operating System :: Microsoft :: Windows',
        'Operating System :: POSIX :: Linux'
        'Development Status :: 1 - Planning',
        'Intended Audience :: Developers',
        'Intended Audience :: Science/Research',
        'License :: Other/Proprietary License',
        'Topic :: Software Development :: Libraries'
        'Topic :: Software Development :: Libraries :: Python Modules',
        'Topic :: Software Development :: Embedded Systems',
        'Topic :: Scientific/Engineering :: Human Machine Interfaces'
    ],
    author='Manish Meganathan, Mariyam A.Ghani',
    author_email='contact.fyrwatch@gmail.com',
    long_description="""
    ========================
    FyrMesh Platform and API
    ========================
    This package contains tools to orchestrate the a mesh of sensor nodes that have 
    a FyrNode Firmware on them, the mesh messaging is built on top of the painlessMesh
    library and implements a custom message object library for the same.

    The package also runs a long running script in the form of a Bottle REST App, 
    which exposes endpoints on the localhost to interact with the mesh. This long 
    running script is responsible for orchestrating the mesh and interacting with
    Firebase for the cloud services.

    The server can also be interacted with by using the CLI tool 'fyrmesh'.
    """
)