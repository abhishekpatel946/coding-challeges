import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.authApp import authApp
from middlewares import Middlewares

app = FastAPI(prefix="/api/v1", docs_url="/docs", openapi_url="/openapi.json",
              redoc_url="/redoc")

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
authApp.add_middleware(Middlewares)

app.mount("/authApp", authApp)


@app.get("/")
async def root():
    return {"status": "true"}


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000, reload=True)
