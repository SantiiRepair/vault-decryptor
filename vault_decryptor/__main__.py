import shutup
import asyncio
import argparse
from pathlib import Path

# used to hide asyncio annoying warning
shutup.please()


async def _main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-ptl",
        "--path-to-logs",
        "--path_to_logs",
        type=str,
        required=True,
        help="Path to metamask log files",
    )
    parser.add_argument(
        "-r",
        "--recursive",
        type=bool,
        default=False,
        help="Iterate over all files in the specified path",
    )

    args = parser.parse_args()
    if args.recursive():
        ptl = Path(args.path_to_logs)
        for logs in ptl.glob("*.log"):
            print(logs)


if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(_main())