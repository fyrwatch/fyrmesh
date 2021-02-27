from setuptools import setup

setup(
    name="fyrmesh",
    version="0.0.1",
    py_modules = ['app'],
    install_requires=[
        'Click'
    ],
    entry_points="""
        [console_scripts]
        fyrmesh=cli:cli
    """
)