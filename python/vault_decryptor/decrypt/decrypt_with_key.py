import json
import base64
from typing import Dict
from termcolor import colored
from Crypto.Cipher import AES


def decrypt_with_key(key: bytes, payload: Dict) -> Dict:
    try:
        data = base64.b64decode(payload["data"])
        iv = base64.b64decode(payload["iv"])
        cipher = AES.new(key, mode=AES.MODE_GCM, nonce=iv)
        decrypted_bytes = cipher.decrypt(data)
        decrypted_text = decrypted_bytes.rstrip(b"\\").decode(
            "utf-8", "ignore"
        )
        print(decrypted_text)
        decrypted = json.loads(decrypted_text)
        return decrypted
    except (Exception, ValueError) as e:
        # if isinstance(e, UnicodeDecodeError):
        # exit(print(colored("[ERROR]: Incorrect Password", "red")))
        print(colored(f"[ERROR]: {e}", "red"))
