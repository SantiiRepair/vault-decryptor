import json
from typing import Any
from termcolor import colored
from vault_decryptor.modules.decrypt_with_key import decrypt_with_key
from vault_decryptor.modules.key_from_password import key_from_password


def decrypt(password: str, text: str, key_bytes: str = None) -> Any:
    try:
        payload = json.loads(text)
    except ValueError as e:
        print(colored(f"[ERROR]: {e}", "red"))
    if "KeyringController" in text:
        salt = payload["KeyringController"]["vault"]["salt"]
        key = key_bytes or key_from_password(password, salt=salt)
        result = decrypt_with_key(key, payload["KeyringController"]["vault"])
        return result
    if "vault" in text:
        salt = payload["vault"]["salt"]
        key = key_bytes or key_from_password(password, salt=salt)
        result = decrypt_with_key(key, payload["vault"])
        return result
    salt = payload["salt"]
    key = key_bytes or key_from_password(password, salt=salt)
    result = decrypt_with_key(key, payload)
    return result
