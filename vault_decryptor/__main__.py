import os
import js2py
import shutup
import asyncio
import argparse
import requests
from pathlib import Path
from decryptor.helpers.bash import bash

# used to hide asyncio annoying warning
shutup.please()
c = os.path.dirname(os.path.abspath(__file__))


async def _main():
    js_packages = ["@metamask/browser-passworder"]
    for package in js_packages:
        await bash(f"npm i -g {package}")


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(_main())
