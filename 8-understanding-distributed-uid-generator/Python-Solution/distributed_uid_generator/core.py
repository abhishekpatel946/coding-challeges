import time

from constants import BinaryBits
from context import get_sequence_number
from utils import (
    get_datacenter_id,
    get_machine_id,
    get_n_binary_bits,
    convert_to_decimal,
)


class DistributedUIDGenerator:
    def __init__(
        self,
        sign_bit=None,
        timestamp=None,
        datacenter_id=None,
        machine_id=None,
        sequence_number=None,
    ):
        self.sign_bit = sign_bit
        self.timestamp = timestamp
        self.datacenter_id = datacenter_id
        self.machine_id = machine_id
        self.sequence_number = sequence_number

    def generate_uid(self):
        sign_bit = get_n_binary_bits(0, BinaryBits.SIGN_BITS_LONG.value)
        timestamp = get_n_binary_bits(
            int(time.time()), BinaryBits.TIMESTAMP_BITS_LONG.value
        )

        sequence_number = get_sequence_number()

        current_datacentar_id = get_datacenter_id()
        datacentar_id = get_n_binary_bits(
            current_datacentar_id, BinaryBits.DATACENTER_ID_TOTAL_LONG.value
        )

        current_machine_id = get_machine_id()
        machine_id = get_n_binary_bits(
            current_machine_id, BinaryBits.MACHINE_ID_TOTAL_LONG.value
        )

        self.__init__(sign_bit, timestamp, datacentar_id, machine_id, sequence_number)

        uid_in_binary = (
            self.sign_bit
            + "-"
            + self.timestamp
            + "-"
            + self.datacenter_id
            + "-"
            + self.machine_id
            + "-"
            + self.sequence_number
        )

        uid_in_number = convert_to_decimal(uid_in_binary)

        return (uid_in_binary, uid_in_number)
