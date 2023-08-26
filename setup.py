#!/usr/bin/env python3
from pathlib import Path

from setuptools import find_packages, setup

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
    url="https://github.com/SantiiRepair/vault-decryptor",
    author="Santiago Ramirez",
    author_email="None",
    license="MIT",
    packages=find_packages(),
    entry_points={
        "console_scripts": ["vault-decryptor = vault_decryptor.__main__:main"]
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
