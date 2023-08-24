import re
import json
from typing import List


# Deduplicates array with rudimentary non-recursive shallow comparison of keys
def dedupe(list: List):
    result = []
    for x in list:
        if not any(
            len(x.keys()) == len(y.keys())
            and all(y[k] == ex for k, ex in x.items())
            for y in result
        ):
            result.append(x)
    return result


def decode_mnemonic(mnemonic):
    if isinstance(mnemonic, str):
        return mnemonic
    else:
        return mnemonic.decode("utf8")


def extract_vault_from_file(data):
    try:
        # Attempt 1: raw JSON
        return json.loads(data)
    except ValueError:
        pass

    # Attempt 2: pre-v3 cleartext
    matches = re.search(r'{"wallet-seed":"([^"}]*)"', data)
    if matches:
        mnemonic = matches.group(1).replace("\\n", '')
        vault_matches = re.search(
            r'"wallet":("{[ -~]*\\"version\\":2}")', data
        )
        vault = json.loads(vault_matches.group(1)) if vault_matches else {}
        return {"data": {"mnemonic": mnemonic, **vault}}

    # Attempt 3: chromium 000003.log file on Linux
    matches = re.search(r'"KeyringController":{"vault":"{[^{}]*}"', data)
    if matches:
        vault_body = matches.group(0)[29:]
        return json.loads(json.loads(vault_body))

    # Attempt 4: chromium 000005.ldb on Windows
    match_regex = r"Keyring[0-9][^\}]*(\{[^\{\}]*\\\"\})"
    capture_regex = r"Keyring[0-9][^\}]*(\{[^\{\}]*\\\"\})"
    iv_regex = r'\\"iv.{1,4}[^A-Za-z0-9+\/]{1,10}([A-Za-z0-9+\/]{10,40}=*)'
    data_regex = r'\\"[^":,is]*\\":\\"([A-Za-z0-9+\/]*=*)'
    salt_regex = (
        r',\\"salt.{1,4}[^A-Za-z0-9+\/]{1,10}([A-Za-z0-9+\/]{10,100}=*)'
    )

    matches = re.findall(match_regex, data)
    vaults = []

    for match in matches:
        capture_match = re.search(capture_regex, match)
        if capture_match:
            data_match = re.search(data_regex, capture_match.group(1))
            iv_match = re.search(iv_regex, capture_match.group(1))
            salt_match = re.search(salt_regex, capture_match.group(1))

            if data_match and iv_match and salt_match:
                vaults.append(
                    {
                        "data": data_match.group(1),
                        "iv": iv_match.group(1),
                        "salt": salt_match.group(1),
                    }
                )

    if not vaults:
        return None

    if len(vaults) > 1:
        print("Found multiple vaults!", vaults)

    return vaults[0]


def is_vault_valid(vault):
    return isinstance(vault, dict) and all(
        isinstance(vault[e], str) for e in ["data", "iv", "salt"]
    )


def decrypt_vault(password, vault):
    if "data" in vault and "mnemonic" in vault["data"]:
        return [vault]

    decrypted_data = passworder.decrypt(password, json.dumps(vault))

    keyrings_with_decoded_mnemonic = []
    for keyring in decrypted_data:
        if "mnemonic" in keyring["data"]:
            keyring["data"]["mnemonic"] = decode_mnemonic(
                keyring["data"]["mnemonic"]
            )
        keyrings_with_decoded_mnemonic.append(keyring)

    return keyrings_with_decoded_mnemonic
