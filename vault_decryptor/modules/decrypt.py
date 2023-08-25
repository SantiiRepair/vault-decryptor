import json
from typing import Any, Optional
from vault_decryptor.modules.decrypt_with_key import decrypt_with_key
from vault_decryptor.modules.key_from_password import key_from_password


async def decrypt(password: str, text: str, key: str = None) -> Any:
    payload = json.loads(text)
    salt = payload["salt"]

    crypto_key = key or await key_from_password(password, salt)

    result = await decrypt_with_key(crypto_key, payload)
    return result
