import shutup
import asyncio
import argparse
from pathlib import Path
from termcolor import colored

# used to hide asyncio annoying warning
shutup.please()


async def _main() -> None:
    parser = argparse.ArgumentParser(
        usage="",
        description=f"{colored('Vault Decryptor is a cli tool that allows you to decrypt vault data of Metamask Extension, this work by entering vault data path and password of the wallet extension, then if the data entered in the arguments are correct it creates a csv file with the seed phrases of the wallet', 'green')}",
    )
    parser.add_argument(
        "-lp",
        "--logs-path",
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

    args = parser.parse_args()
    if args.recursive == "yes":
        ptl = Path(args.logs_path)
        for logs in ptl.glob("*.log"):
            print(logs)


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(_main())
