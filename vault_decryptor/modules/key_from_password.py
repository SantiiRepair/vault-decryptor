from cryptography.hazmat.primitives.hashes import SHA256
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC


def key_from_password(password: str, salt: str):
    password_bytes = password.encode("utf-8")
    salt_bytes = salt.encode("utf-8")

    kdf = PBKDF2HMAC(
        algorithm=SHA256(),
        length=32,  # 256 bits key length
        salt=salt_bytes,
        iterations=10000,
    )

    return kdf.derive(password_bytes)
