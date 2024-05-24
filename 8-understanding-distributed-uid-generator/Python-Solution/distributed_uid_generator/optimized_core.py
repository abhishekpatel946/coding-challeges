import time
import threading

# Constants for Snowflake
EPOCH = int(
    time.time() * 1000
)  # Custom epoch timestamp in milliseconds (used to reduce the size of the timestamp part)
DATACENTER_ID_BITS = 5
MACHINE_ID_BITS = 5
SEQUENCE_BITS = 12

MAX_DATACENTER_ID = -1 ^ (-1 << DATACENTER_ID_BITS)
MAX_MACHINE_ID = -1 ^ (-1 << MACHINE_ID_BITS)
MAX_SEQUENCE = -1 ^ (-1 << SEQUENCE_BITS)

# Shifts (Bit shifts to place each component in the correct position in the final ID)
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
        self.lock = threading.Lock()  # "lock" for thread safety

    # Returns the current timestamp in milliseconds
    def _time_gen(self):
        return int(time.time() * 1000)

    # Waits until the next millisecond if the clock has not moved forward
    def _wait_for_next_millis(self, last_timestamp):
        timestamp = self._time_gen()
        while timestamp <= last_timestamp:
            timestamp = self._time_gen()
        return timestamp

    def generate_id(self):
        # Acquires a lock for thread safety
        with self.lock:
            # Gets the current timestamp
            timestamp = self._time_gen()

            # Checks if the clock has moved backwards and raises an exception if it has
            if timestamp < self.last_timestamp:
                raise Exception("Clock moved backwards. Refusing to generate id")

            # If the timestamp is the same as the last one, increments the sequence number.
            # If the sequence number overflows, waits for the next millisecond
            if timestamp == self.last_timestamp:
                self.sequence = (self.sequence + 1) & MAX_SEQUENCE
                if self.sequence == 0:
                    timestamp = self._wait_for_next_millis(self.last_timestamp)
            else:
                # If the timestamp is different, resets the sequence number
                self.sequence = 0

            # Updates the last timestamp
            self.last_timestamp = timestamp

            # Constructs the ID by combining the timestamp, datacenter ID, machine ID, and sequence number with bitwise operations
            id = (
                ((timestamp - EPOCH) << TIMESTAMP_SHIFT)
                | (self.datacenter_id << DATACENTER_ID_SHIFT)
                | (self.machine_id << MACHINE_ID_SHIFT)
                | self.sequence
            )

            # Returns the generated ID
            return id


class TestSnowflakeGenerator:

    def __init__(self, total_num_ids):
        self.generator = SnowflakeGenerator(datacenter_id=1, machine_id=1)
        self.num_ids = total_num_ids

    def test_high_throughput(self):
        ids = set()

        for _ in range(self.num_ids):
            new_id = self.generator.generate_id()
            if new_id in ids:
                print("Duplicate ID found:", new_id)
                return False
            ids.add(new_id)

        return len(ids) == self.num_ids

    def test(self):
        start_time = time.time()
        if self.test_high_throughput():
            duration = time.time() - start_time
            print(f"Successfully generated unique IDs, Generated {self.num_ids} IDs in {duration:.3f} seconds")
            return duration
        else:
            print("Failed to generate unique IDs.")

