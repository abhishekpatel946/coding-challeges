import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.authApp import authApp
from middlewares import RateLimiterMiddleware

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
authApp.add_middleware(RateLimiterMiddleware)

app.mount("/authApp", authApp)


@app.get("/")
async def root():
    return {"status": "true"}


# @app.get("/unlimited")
# async def unlimited():
#     return {"message": "Unlimited! Let's Go!"}


# @app.get("/limited", dependencies=[Depends(RateLimiter())])
# async def limited():
#     return {"message": "Limited, don't over use me!"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000, reload=True)
