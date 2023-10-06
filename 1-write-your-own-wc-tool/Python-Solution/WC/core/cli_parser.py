import argparse

from WC import __version__


class WCArgumentParser:
    def __init__(self):

        self.parser = argparse.ArgumentParser(
            epilog="To get addtional help on a specific provider run: python wc.py -h")

        self.parser.add_argument(
            "-v", "--version", action="version", version=f"WC {__version__}")

        self._init_common_args_parser()

    def _init_common_args_parser(self):
        parser = self.parser.add_argument_group("WC Arguments")

        parser.add_argument(
            "-f",
            "--filename",
            dest="filename",
            default=None,
            # type=str,
            help="The input file, or standard input (if no file is specified) to the standard output.",
        )

        parser.add_argument(
            "-c",
            "--bytes",
            dest="bytes",
            default=False,
            action="store_true",
            help="The number of bytes in each input file is written to the standard output.",
        )
        parser.add_argument(
            "-l",
            "--lines",
            dest="lines",
            default=False,
            action="store_true",
            help="The number of lines in each input file is written to the standard output.",
        )
        parser.add_argument(
            "-m",
            "--characters",
            dest="chars",
            default=False,
            action="store_true",
            help="The number of characters in each input file is written to the standard output.",
        )
        parser.add_argument(
            "-w",
            "--words",
            dest="words",
            default=False,
            action="store_true",
            help="The number of words in each input file is written to the standard output.",
        )

    def parse_args(self, args=None):
        args = self.parser.parse_args(args)
        return args
