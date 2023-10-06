#!/usr/bin/env python
import os

from setuptools import find_packages
from setuptools import setup

# Package info
NAME = "WC Tool"
ROOT = os.path.dirname(__file__)
VERSION = __import__(NAME).__version__

# Requirements
requirements = []
with open(
    os.path.join(os.path.dirname(os.path.realpath(__file__)), "requirements.txt")
) as f:
    for r in f.readlines():
        requirements.append(r.strip())

# Setup
setup(
    name=NAME,
    version=VERSION,
    description="WC – word, line, character, and byte count",
    long_description_content_type="text/markdown",
    long_description=open("README.md").read(),
    author="Abhishek Patel",
    url="https://github.com/abhishekpatel946/",
    packages=find_packages(),
    include_package_data=True,
    install_requires=requirements,
    license="GNU General Public License v2 (GPLv2)",
    classifiers=[
        "Intended Audience :: Developers",
        "Natural Language :: English",
        "License :: OSI Approved :: GNU General Public License v2 (GPLv2)",
        "Programming Language :: Python",
    ],
)
