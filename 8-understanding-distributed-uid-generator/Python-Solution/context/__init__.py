import asyncio
import contextlib

from utils import get_n_binary_bits
from constants import BinaryBits

# global state variables
_sequence_number = 0
MAX_LIMIT_FOR_SEQUENCE_NUMBER = 999


@contextlib.asynccontextmanager
async def increment_sequence(bits: int = 1, every_ms: int = 1):
    global _sequence_number
    try:
        while True:
            if _sequence_number > MAX_LIMIT_FOR_SEQUENCE_NUMBER:
                _sequence_number = 0
            _sequence_number += bits
            await asyncio.sleep(every_ms / 1000)  # Convert milliseconds to seconds
            # print("sequence number", _sequence_number)
    except asyncio.CancelledError as e:
        print(f"Sequence increment cancelled: {e}")
    finally:
        yield


async def run_increment_sequence():
    async with increment_sequence(bits=1, every_ms=1):
        await asyncio.sleep(1)  # Run for 1 seconds


# @contextlib.contextmanager
def get_sequence_number():
    # declare the "_sequence_number" that we're using below
    global _sequence_number

    # return the sequence number in binary format
    seq_number = get_n_binary_bits(_sequence_number, BinaryBits.SEQUANCE_NUMBER_BITS_LONG.value)
    # print(f"sequence number {_sequence_number} in binary format {seq_number}")

    # return
    return seq_number
