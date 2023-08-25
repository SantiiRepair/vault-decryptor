import logging


def logger(boolean):
    logging.basicConfig(
        format="[%(levelname) 5s/%(asctime)s] %(name)s: %(message)s",
        level=logging.DEBUG if boolean else logging.INFO,
    )
