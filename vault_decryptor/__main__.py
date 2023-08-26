import json
import shutup
import logging
import argparse
from pathlib import Path
from termcolor import colored
from vault_decryptor.modules.decrypt import decrypt

# used to hide asyncio annoying warning
shutup.please()


def main() -> None:
    # loop = asyncio.get_event_loop()
    # loop.run_until_complete(main())
    parser = argparse.ArgumentParser(
        prog=colored("vault-decryptor", "yellow"),
        usage=colored("vault-decryptor [-r] [-l]", "green"),
        epilog=colored("Thanks for use %(prog)s!", "green"),
        description=colored(
            "Vault Decryptor is a cli tool that allows you to decrypt vault data of Metamask Extension, this work by entering vault data path and password of the wallet extension, then if the data entered in the arguments are correct it creates a csv file with the seed phrases of the wallet",
            "green",
        ),
    )
    parser.add_argument(
        "-l",
        "--log",
        type=str,
        help="Path to metamask log file",
    )

    parser.add_argument(
        "-p",
        "--password",
        type=str,
        help="Password of your metamask vault",
    )

    parser.add_argument(
        "-r",
        "--recursive",
        type=str,
        default="no",
        help="Iterate over all files in the specified path",
    )

    parser.add_argument(
        "-v",
        "--vault",
        type=str,
        help="Path to vault json file",
    )

    parser.add_argument(
        "-d",
        "--debug",
        type=str,
        default="no",
        help="Enable logging framework debug mode",
    )

    args = parser.parse_args()
    if args.debug == "yes":
        logger(True)
    logger(False)
    if args.recursive == "yes":
        if args.log:
            logs = Path(args.log)
            for log in logs.glob("*.log"):
                print(log)
        if args.vault:
            if not args.password:
                exit(
                    print(
                        colored(
                            "[ERROR]: Metamask Password is required", "red"
                        )
                    )
                )
            if ".json" in args.vault:
                with open(args.vault, "r") as jsf:
                    payload = jsf.read()
                    result = decrypt(password=args.password, text=payload)
                    print(result)
                    return result
            vaults = Path(args.vault)
            for vault in vaults.glob("*.json"):
                with open(vault, "r") as v:
                    vault_data = json.load(v)
                    print(vault_data)


def logger(boolean: bool):
    logging.basicConfig(
        format="[%(levelname) 5s/%(asctime)s] %(name)s: %(message)s",
        level=logging.DEBUG if boolean else logging.INFO,
    )
