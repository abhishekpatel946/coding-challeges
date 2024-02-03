from fastapi import FastAPI

authApp = FastAPI(prefix="/api/v1", docs_url="/docs", openapi_url="/openapi.json",
                  redoc_url="/redoc")


@authApp.get("/limited")
async def limited():
    return {"message": "Limited, don't over use me!"}
