import base64
from Crypto.Hash import SHA512
from Crypto.Protocol.KDF import PBKDF2


def key_from_password(password: str, salt: str):
    password_bytes = password.encode("utf-8")
    salt_bytes = base64.b64decode(salt)

    key = PBKDF2(
        password_bytes,
        salt_bytes,
        dkLen=32,
        count=10000,
        hmac_hash_module=SHA512,
    )

    return key
