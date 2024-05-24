import time
import threading

from main import SnowflakeGenerator


# Example test for generating 10,000 IDs in one second
def test_high_throughput(generator, num_ids):
    ids = set()
    start_time = time.time()

    for _ in range(num_ids):
        new_id = generator.generate_id()
        if new_id in ids:
            print("Duplicate ID found:", new_id)
            return False
        ids.add(new_id)

    duration = time.time() - start_time
    print(f"Generated {num_ids} IDs in {duration:.3f} seconds")

    return len(ids) == num_ids


if __name__ == "__main__":
    generator = SnowflakeGenerator(datacenter_id=1, machine_id=1)
    num_ids = 10000

    if test_high_throughput(generator, num_ids):
        print("Successfully generated unique IDs.")
    else:
        print("Failed to generate unique IDs.")
