import uuid
from fastapi import Request
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.responses import JSONResponse
from fastapi import status as HttpStatus


class RequestIDMiddleware(BaseHTTPMiddleware):
    def __init__(self, app):
        super().__init__(app)

    async def dispatch(self, request: Request, call_next):
        try:
            request.state.request_id = str(uuid.uuid4())
            response = await call_next(request)
            response.headers["X-Request-ID"] = request.state.request_id
            return response

        except Exception as e:
            return JSONResponse(
                content={
                    "responseData": {},
                    "message": [{"msg": f"Something went wrong {e}"}],
                    "success": False,
                    "code": HttpStatus.HTTP_500_INTERNAL_SERVER_ERROR,
                },
                status_code=HttpStatus.HTTP_500_INTERNAL_SERVER_ERROR,
            )
