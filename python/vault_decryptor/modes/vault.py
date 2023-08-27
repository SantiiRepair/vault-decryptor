import sys
from pathlib import Path
from termcolor import colored
from vault_decryptor.modules.decrypt import decrypt


def vault(path: str, password: str, key_bytes="", recursive="no"):
    if recursive == "yes":
        if not password:
            sys.exit(
                print(colored("[ERROR]: Metamask Password is required", "red"))
            )
        if ".json" in path:
            sys.exit(
                print(
                    colored(
                        "[ERROR]: Recursive mode expect folder path, not file",
                        "red",
                    )
                )
            )
        vaults = Path(path)
        for vault_path in vaults.glob("*.json"):
            with open(vault_path, "r") as jsf:
                text = jsf.read()
                decrypt(password, text, key_bytes)
    with open(path, "r") as jsf:
        text = jsf.read()
        decrypt(password, text, key_bytes)
