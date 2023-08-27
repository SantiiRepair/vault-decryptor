import json
import base64
from typing import Dict
from termcolor import colored
from Crypto.Cipher import AES


def decrypt_with_key(key: bytes, payload: Dict) -> Dict:
    try:
        data = base64.b64decode(payload["data"])[:-16]
        iv = base64.b64decode(payload["iv"])
        cipher = AES.new(key, mode=AES.MODE_GCM, nonce=iv)
        decrypted_data = cipher.decrypt(data)
        print(decrypted_data)
        decrypted = json.loads(decrypted_data)
        return decrypted
    except (Exception, ValueError) as e:
        # if isinstance(e, UnicodeDecodeError):
        # exit(print(colored("[ERROR]: Incorrect Password", "red")))
        print(colored(f"[ERROR]: {e}", "red"))
