import json
import os
import queue
import sys
import time
from threading import Thread

import requests

NO_OF_REQUESTS_PER_SECOND = 10_000  # requests per second
URL = "http://localhost:8000/authApp/optimized-generate"
CONCURRENT = 1

queue = queue.Queue(CONCURRENT * 2)
responses = []


# utils
def write_file_content(file_path=None, content=None) -> bool:
    with open(file_path, "w+") as f:
        json.dump(content, f)


def create_directory(dir_path=None) -> bool:
    if not os.path.exists(dir_path):
        os.makedirs(dir_path)
        return True
    return False


def make_request():
    while True:
        url = queue.get()
        try:
            response = requests.get(url)
            if response.status_code == 200:
                responses.append(json.loads(response.text))
        except requests.RequestException as e:
            print(f"Request failed: {e}")
        finally:
            queue.task_done()


def run_unittests():
    print(
        f"run-unittest... \n total no of request {NO_OF_REQUESTS_PER_SECOND} \n with concurrency {CONCURRENT}"
    )

    for i in range(CONCURRENT):
        t = Thread(target=make_request)
        t.setDaemon(True)
        t.start()

    try:
        for _ in range(NO_OF_REQUESTS_PER_SECOND):
            queue.put(URL)
        queue.join()
    except KeyboardInterrupt:
        sys.exit(1)
    finally:
        lst = list()
        duplicate_count = 0
        for item in responses:
            uid = item["responseData"]["unique_id"]
            if uid not in lst:
                lst.append(item["responseData"]["unique_id"])
            else:
                duplicate_count += 1
        print("total duplicates are:", duplicate_count, "out of", len(responses))

        # try:
        #     if responses and len(responses):
        #         # genereate the "result/*" directory
        #         cwd = os.getcwd()
        #         result_dir = f"{cwd}/result"
        #         print(f"Storing the outcomes into {result_dir}")
        #         create_directory(dir_path=result_dir)
        #         write_file_content(
        #             file_path=f"{result_dir}/test-result.json",
        #             content=responses,
        #         )
        #         print(
        #             f"Successfully stored outcome & remove the ouput directory {result_dir}/result.json"
        #         )
        # except Exception as e:
        #     print("Failed to store the changes in result directory", e)
        # finally:
        #     print("Exiting test...\t Goodbye!!!")


if __name__ == "__main__":
    start_time = time.time()

    # run tests
    run_unittests()

    # print the time elapsed
    duration = time.time() - start_time
    print(
        f"Successfully generated unique IDs, Generated {NO_OF_REQUESTS_PER_SECOND} IDs in {duration:.3f} seconds"
    )
