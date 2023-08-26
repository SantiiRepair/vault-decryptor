import base64
import hashlib
from Crypto.Protocol.KDF import PBKDF2


def key_from_password(password, salt):
    pass_buffer = password.encode("utf-8")
    salt_buffer = base64.b64decode(salt)

    key = PBKDF2(
        pass_buffer,
        salt_buffer,
        dkLen=32,
        count=10000,
        hmac_hash_module=hashlib.sha256,
    )

    return key
