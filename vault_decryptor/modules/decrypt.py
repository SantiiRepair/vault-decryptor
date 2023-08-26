import json
from typing import Any
from termcolor import colored
from vault_decryptor.modules.decrypt_with_key import decrypt_with_key
from vault_decryptor.modules.key_from_password import key_from_password


def decrypt(password: str, text: str, key: str = None) -> Any:
    try:
        payload = json.loads(text)
        if "KeyringController" in text:
            salt = payload["KeyringController"]["vault"]["salt"]
            crypto_key = key or key_from_password(password, salt)
            result = decrypt_with_key(crypto_key, payload)
            return result
        if "vault" in text:
            salt = payload["vault"]["salt"]
            crypto_key = key or key_from_password(password, salt)
            result = decrypt_with_key(crypto_key, payload)
            return result
        salt = payload["salt"]
        crypto_key = key or key_from_password(password, salt)
        result = decrypt_with_key(crypto_key, payload)
        return result
    except ValueError:
        print(colored("[ERROR]: The JSON file is invalid", "red"))
