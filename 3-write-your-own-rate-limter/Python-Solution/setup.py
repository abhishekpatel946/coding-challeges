#!/usr/bin/env python
import os

from setuptools import find_packages
from setuptools import setup

# Package info
NAME = "Rate Limter"
ROOT = os.path.dirname(__file__)
VERSION = __import__(NAME).__version__

# Requirements
requirements = []
with open(
    os.path.join(os.path.dirname(
        os.path.realpath(__file__)), "requirements.txt")
) as f:
    for r in f.readlines():
        requirements.append(r.strip())

# Setup
setup(
    name=NAME,
    version=VERSION,
    description="Rate Limter – limit the api requests & protect against concurrent requests like ddos attacks",
    long_description_content_type="text/markdown",
    long_description=open("README.md").read(),
    author="Abhishek Patel",
    url="https://github.com/abhishekpatel946/",
    packages=find_packages(),
    include_package_data=True,
    install_requires=requirements,
    license="MIT",
    classifiers=[
        "Intended Audience :: Developers",
        "Natural Language :: English",
        "License :: OSI Approved :: MIT",
        "Programming Language :: Python",
    ],
)
