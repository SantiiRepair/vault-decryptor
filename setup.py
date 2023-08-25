#!/usr/bin/env python3
from pathlib import Path

import setuptools
from setuptools import setup

this_dir = Path(__file__).parent
module_dir = this_dir / "vault_decryptor"

requirements = []
requirements_path = this_dir / "requirements.txt"
if requirements_path.is_file():
    with open(requirements_path, "r", encoding="utf-8") as requirements_file:
        requirements = requirements_file.read().splitlines()

# -----------------------------------------------------------------------------

setup(
    name="vault-decryptor",
    version="1.0.0",
    description="A fast, local Metamask Vault Decryptor in the command line.",
    url="http://github.com/SantiiRepair/vault-decryptor",
    author="Santiago Ramirez",
    license="MIT",
    packages=setuptools.find_packages(),
    entry_points={
        "console_scripts": [
            "vault_decryptor = vault_decryptor.__main__:main",
        ]
    },
    install_requires=requirements,
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
    ],
    keywords="SantiiRepair vault-decryptor",
)
