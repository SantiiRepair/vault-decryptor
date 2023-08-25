import base64
import hashlib
from Crypto.Protocol.KDF import PBKDF2
from Crypto.Cipher import AES

STRING_ENCODING = "utf-8"


async def key_from_password(password, salt, exportable=False):
    pass_buffer = password.encode(STRING_ENCODING)
    salt_buffer = base64.b64decode(salt)

    key = PBKDF2(
        pass_buffer,
        salt_buffer,
        dkLen=32,
        count=10000,
        hmac_hash_module=hashlib.sha256,
    )

    return key
