from fastapi import FastAPI
from fastapi import status as HttpStatus
from pydantic import BaseModel
from starlette.responses import JSONResponse

from context import _global_state

authApp = FastAPI(prefix="/api/v1", docs_url="/docs", openapi_url="/openapi.json",
                  redoc_url="/redoc")


class RequestBody(BaseModel):
    user_id: int = 1
    user_name: str = "johndoe"
    password: str = "password123"


@authApp.get("/limited")
async def limited():
    return JSONResponse(
        content={
            "responseData": {},
            "message": f"Limited requests, Do not overuse me. Maximum request {_global_state['user_state']['capacity']} is allowed.",
            "success": False,
            "code": HttpStatus.HTTP_429_TOO_MANY_REQUESTS,
        },
        status_code=HttpStatus.HTTP_429_TOO_MANY_REQUESTS,
    )
