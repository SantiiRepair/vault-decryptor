import json
import base64
from typing import Any
from termcolor import colored
from Crypto.Cipher import AES


def decrypt_with_key(key: str, payload: dict) -> Any:
    encrypted_data = base64.b64decode(payload["data"])
    vector = base64.b64decode(payload["iv"])

    try:
        cipher = AES.new(key=key, mode=AES.MODE_GCM, iv=vector)
        decrypted_data = cipher.decrypt(encrypted_data)
        decrypted_str = decrypted_data.decode("utf-8").rstrip("\0")
        decrypted_obj = json.loads(decrypted_str)
        return decrypted_obj
    except ValueError as e:
        print(colored(f"[ERROR]: {e}", "red"))
