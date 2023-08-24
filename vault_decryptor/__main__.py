import os
import js2py
import shutup
import asyncio
import argparse
import requests
from pathlib import Path
from decryptor.helpers.bash import bash

shutup.please()
c = os.path.dirname(os.path.abspath(__file__))


async def _main():
    js_packages = ["@metamask/browser-passworder"]
    jsmod_dir = Path(f"{c}/jsmod")
    pymod_dir = Path(f"{c}/pymod")
    passworder_raw = "https://raw.githubusercontent.com/MetaMask/vault-decryptor/master/app/lib.js"
    if not os.path.exists(jsmod_dir):
        os.makedirs(jsmod_dir)
    if not os.path.exists(pymod_dir):
        os.makedirs(pymod_dir)
    if not os.path.exists(f"{jsmod_dir}/passworder.js"):
        passworderjs = requests.get(url=passworder_raw)
        with open(f"{jsmod_dir}/passworder.js", "w") as f:
            f.write(passworderjs.content.decode())
    for package in js_packages:
        await bash(f"npm i -g {package}")
    for jsmod in jsmod_dir.glob("*.js"):
        js2py.translate_file(
            str(jsmod),
            str(jsmod).replace("jsmod", "pymod").replace(".js", ".py"),
        )


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(_main())
