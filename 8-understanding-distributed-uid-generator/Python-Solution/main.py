import asyncio
import threading
import psutil

import uvicorn
from app.authApp import authApp
from context import run_increment_sequence
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from middlewares import RequestIDMiddleware

cpu_cores = psutil.cpu_count(logical=True)
num_workers = cpu_cores * 2 + 1

app = FastAPI(
    prefix="/api/v1", docs_url="/docs", openapi_url="/openapi.json", redoc_url="/redoc"
)

origins = ["*"]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
authApp.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
authApp.add_middleware(RequestIDMiddleware)
app.mount("/authApp", authApp)


@app.get("/")
async def root():
    return {"status": "true"}


# run the async function in event loop in background
def run_async_in_thread(async_func):
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    loop.run_until_complete(async_func())
    loop.close()


# initialize the fastAPI server in main thread
def init_fastAPI():
    uvicorn.run("main:app", host="0.0.0.0", port=8000, workers=2)


if __name__ == "__main__":
    # create a new thread to run incr_sequence
    t = threading.Thread(target=run_async_in_thread, args=(run_increment_sequence,))
    t.daemon = (
        True  # set the thread as a daemon so it will exit when the main thread exits
    )
    t.start()

    # # Start the FastAPI server in the main thread
    uvicorn.run(
        "main:app", host="0.0.0.0", port=8000, workers=1
    )  # run on single thread

    # Optionally, you can join the thread if you want the main thread to wait for it
    # t.join()
