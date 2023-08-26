import json
import base64
from typing import Dict
from termcolor import colored
from cryptography.hazmat.primitives.ciphers.aead import AESGCM


def decrypt_with_key(key: bytes, payload: Dict) -> Dict:
    try:
        iv = base64.b64decode(payload["iv"])
        encrypted = base64.b64decode(payload["data"])
        data = encrypted[0 : len(encrypted) - 16]
        tag = encrypted[len(encrypted) - 16 : len(encrypted)]
        decrypted_data = AESGCM(key).decrypt(
            iv, data=data, associated_data=tag
        )
        print(decrypted_data)
        decrypted = json.loads(decrypted_data)
        print(decrypted)
        return decrypted
    except ValueError as e:
        print(colored(f"[ERROR]: {e}", "red"))
