
import os

from Parser.core.helper import Helper
from Parser.core.lexer import Lexer


def run_from_cli():
    try:
        json = {"key": "value"}

        tokens = Lexer.analyse(json)
        print(tokens)

    except (KeyboardInterrupt, SystemExit):
        print("Exiting")
        return 130
