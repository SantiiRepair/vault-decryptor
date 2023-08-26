import shutup
import logging
import argparse
from termcolor import colored
from vault_decryptor.modes.vault import vault

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
        "-k",
        "--key",
        type=str,
        help="PBKDF2 derived key if you have any",
    )

    parser.add_argument(
        "-l",
        "--log",
        type=str,
        help="Path to metamask log file",
    )

    parser.add_argument(
        "-m",
        "--mode",
        type=str,
        required=True,
        help="Run tool mode, log or vault ",
    )

    parser.add_argument(
        "-path",
        "--path",
        type=str,
        help="Path to log or vault, folder or file",
    )

    parser.add_argument(
        "-pass",
        "--password",
        type=str,
        help="Password of your metamask vault",
    )

    parser.add_argument(
        "-r",
        "--recursive",
        type=str,
        help="Iterate over all files in the specified path",
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
    if args.mode == "log":
        vault(args.path, args.password, args.key, args.recursive)
    if args.mode == "vault":
        vault(args.path, args.password, args.key, args.recursive)


def logger(boolean: bool):
    logging.basicConfig(
        format="[%(levelname) 5s/%(asctime)s] %(name)s: %(message)s",
        level=logging.DEBUG if boolean else logging.INFO,
    )
