import json
import base64
from typing import Dict
from termcolor import colored
from cryptography.exceptions import InvalidKey, InvalidTag
from cryptography.hazmat.primitives.ciphers.aead import AESGCM


def decrypt_with_key(key: bytes, payload: Dict) -> Dict:
    try:
        encrypted = base64.b64decode(payload["data"])
        iv = base64.b64decode(payload["iv"])
        data = encrypted[0 : len(encrypted) - 16]
        tag = encrypted[len(encrypted) - 16 : len(encrypted)]
        decrypted_data = AESGCM(key).decrypt(
            iv, data=data, associated_data=tag
        )
        print(decrypted_data)
        decrypted = json.loads(decrypted_data)
        print(decrypted)
        return decrypted
    except (Exception, ValueError) as e:
        if isinstance(e, InvalidTag):
            exit(print(colored("[ERROR]: Invalid Tag", "red")))
        if isinstance(e, InvalidKey):
            exit(print(colored("[ERROR]: Invalid Key", "red")))
        print(colored(f"[ERROR]: {e}", "red"))
