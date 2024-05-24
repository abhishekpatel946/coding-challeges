import json
import os
import time

from constants import BinaryBits

# ---------------------------------------------------------------- #
# file I/O utilities


def is_file_exists(file_path=None) -> bool:
    return os.path.exists(file_path)


def write_file_content(file_path=None, content=None) -> bool:
    if not is_file_exists(file_path):
        with open(file_path, "w+") as f:
            json.dump(content, f)
        return True
    print("File path does not exist")
    return False


def create_directory(dir_path=None) -> bool:
    if not os.path.exists(dir_path):
        os.makedirs(dir_path)
        return True
    return False


# ---------------------------------------------------------------- #
# calculation logic utilities


def get_n_binary_bits(number: int, bits: int) -> str:
    binary_str = bin(number)[2:]
    return binary_str.zfill(bits)


def get_datacenter_id():
    return (
        int(time.time()) % BinaryBits.DATACENTER_ID_TOTAL_BITS.value
    )  # datacenter_id is 5 bits long


def get_machine_id():
    return (
        int(time.time()) % BinaryBits.MACHINE_ID_TOTAL_BITS.value
    )  # machine_id is 5 bits long


def convert_to_decimal(value: str) -> int:
    # remove the "-" from string
    value = value.replace('-', '')
    # return decimal formatted value
    return int(value, 2)
