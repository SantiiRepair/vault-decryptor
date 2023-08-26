import shutup
import logging
# import asyncio
import argparse
from pathlib import Path
from termcolor import colored

# used to hide asyncio annoying warning
shutup.please()


def main() -> None:
    # loop = asyncio.get_event_loop()
    # loop.run_until_complete(main())
    parser = argparse.ArgumentParser(
        prog=f"{colored('vault-decryptor', 'yellow')}",
        usage=f"{colored('vault-decryptor [-r] [-l]', 'green')}",
        epilog=f"{colored('Thanks for use %(prog)s!', 'green')}",
        description=f"{colored('Vault Decryptor is a cli tool that allows you to decrypt vault data of Metamask Extension, this work by entering vault data path and password of the wallet extension, then if the data entered in the arguments are correct it creates a csv file with the seed phrases of the wallet', 'green')}",
    )
    parser.add_argument(
        "-l",
        "--log",
        type=str,
        help="Path to metamask log files",
    )

    parser.add_argument(
        "-r",
        "--recursive",
        type=str,
        default="no",
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
    if args.recursive == "yes":
        log = Path(args.log)
        for logs in log.glob("*.log"):
            print(logs)


def logger(boolean):
    logging.basicConfig(
        format="[%(levelname) 5s/%(asctime)s] %(name)s: %(message)s",
        level=logging.DEBUG if boolean else logging.INFO,
    )
