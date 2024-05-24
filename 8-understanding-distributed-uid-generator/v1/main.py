import time
import threading

# Constants for Snowflake
EPOCH = 1288834974657  # Custom epoch timestamp in milliseconds
DATACENTER_ID_BITS = 5
MACHINE_ID_BITS = 5
SEQUENCE_BITS = 12

MAX_DATACENTER_ID = -1 ^ (-1 << DATACENTER_ID_BITS)
MAX_MACHINE_ID = -1 ^ (-1 << MACHINE_ID_BITS)
MAX_SEQUENCE = -1 ^ (-1 << SEQUENCE_BITS)

# Shifts
TIMESTAMP_SHIFT = DATACENTER_ID_BITS + MACHINE_ID_BITS + SEQUENCE_BITS
DATACENTER_ID_SHIFT = MACHINE_ID_BITS + SEQUENCE_BITS
MACHINE_ID_SHIFT = SEQUENCE_BITS

class SnowflakeGenerator:
    def __init__(self, datacenter_id, machine_id):
        if datacenter_id > MAX_DATACENTER_ID or datacenter_id < 0:
            raise ValueError(f"datacenter_id must be between 0 and {MAX_DATACENTER_ID}")
        if machine_id > MAX_MACHINE_ID or machine_id < 0:
            raise ValueError(f"machine_id must be between 0 and {MAX_MACHINE_ID}")
        
        self.datacenter_id = datacenter_id
        self.machine_id = machine_id
        self.sequence = 0
        self.last_timestamp = -1
        self.lock = threading.Lock()

    def _time_gen(self):
        return int(time.time() * 1000)

    def _wait_for_next_millis(self, last_timestamp):
        timestamp = self._time_gen()
        while timestamp <= last_timestamp:
            timestamp = self._time_gen()
        return timestamp

    def generate_id(self):
        with self.lock:
            timestamp = self._time_gen()

            if timestamp < self.last_timestamp:
                raise Exception("Clock moved backwards. Refusing to generate id")

            if timestamp == self.last_timestamp:
                self.sequence = (self.sequence + 1) & MAX_SEQUENCE
                if self.sequence == 0:
                    timestamp = self._wait_for_next_millis(self.last_timestamp)
            else:
                self.sequence = 0

            self.last_timestamp = timestamp

            id = ((timestamp - EPOCH) << TIMESTAMP_SHIFT) | \
                 (self.datacenter_id << DATACENTER_ID_SHIFT) | \
                 (self.machine_id << MACHINE_ID_SHIFT) | \
                 self.sequence

            return id

# Example usage
if __name__ == "__main__":
    generator = SnowflakeGenerator(datacenter_id=1, machine_id=1)

    for _ in range(10):  # Generate 10 unique IDs
        print(generator.generate_id())
