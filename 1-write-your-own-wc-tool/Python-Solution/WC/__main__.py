
import os

from WC.core.cli_parser import WCArgumentParser
from WC.core.helper import Helper


def run_from_cli():
    try:
        parser = WCArgumentParser()
        args = parser.parse_args()

        # Get the dictionary to get None instead of a crash
        args = args.__dict__

        file_name = None
        if args.get('filename') is not None:
            file_name = args.get('filename')

        file_path = os.path.join(os.getcwd(), file_name)

        if file_path is not None and os.path.exists(file_path):
            if args.get('bytes'):
                size_in_bytes = Helper.get_bytes(file_name)
                print("size in bytes:", size_in_bytes)
            elif args.get('lines'):
                total_lines = Helper.get_lines(file_name)
                print('total lines:', total_lines)
            elif args.get('chars'):
                total_characters = Helper.get_characters(file_name)
                print('total characters:', total_characters)
            elif args.get('words'):
                total_words = Helper.get_words(file_path)
                print('total words:', total_words)
            else:
                size_in_bytes = Helper.get_bytes(file_name)
                total_lines = Helper.get_lines(file_name)
                total_characters = Helper.get_characters(file_name)
                total_words = Helper.get_words(file_path)
                print("size in bytes:", size_in_bytes)
                print('total lines:', total_lines)
                print('total characters:', total_characters)
                print('total words:', total_words)
        else:
            print(f"Could not find file {file_name}")

    except (KeyboardInterrupt, SystemExit):
        print("Exiting")
        return 130
