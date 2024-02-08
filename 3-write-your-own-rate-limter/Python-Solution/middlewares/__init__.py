import os

import yaml
from fastapi import Request
from fastapi import status as HttpStatus
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.responses import JSONResponse

from context import update_global_context
from core import RateLimiterStrategyFactory
from mock_database import MockDB


class Middlewares(BaseHTTPMiddleware):
    def __init__(self, app):
        super().__init__(app)

    async def dispatch(self, request: Request, call_next):
        try:
            await CallerIdentityMiddlware.get_caller_identity(request)
            response_code = await RateLimiterMiddleware.set_rate_limiter_configuration(request)

            if response_code == HttpStatus.HTTP_429_TOO_MANY_REQUESTS:
                return JSONResponse(
                    content={
                        "responseData": {},
                        "message": "Too many requests, please try again later.",
                        "success": False,
                        "code": HttpStatus.HTTP_429_TOO_MANY_REQUESTS,
                    },
                    status_code=HttpStatus.HTTP_429_TOO_MANY_REQUESTS,
                )

            # set the request parameters if required
            next_response = await call_next(request)

            return next_response

        except Exception as e:
            return JSONResponse(
                content={
                    "responseData": {},
                    "message": [{"msg": f"Something went wrong {e}"}],
                    "success": False,
                    "code": HttpStatus.HTTP_400_BAD_REQUEST,
                },
                status_code=HttpStatus.HTTP_400_BAD_REQUEST,
            )


class RateLimiterMiddleware(BaseHTTPMiddleware):
    async def set_rate_limiter_configuration(request: Request):
        try:
            cwd = os.getcwd()
            with open(f"{cwd}/rate_limit.config.yml", "r") as f:
                configuration = yaml.safe_load(f)

            rate_limter_strategy = None
            if 'ratelimiter' in configuration:
                if 'algorithm' in configuration['ratelimiter']:
                    algorithms = configuration["ratelimiter"]['algorithm']
                    for algo in algorithms:
                        if algorithms[algo]['active'] == 'enabled':
                            rate_limter_strategy = RateLimiterStrategyFactory.get_strategy(
                                algo)
                            break

            if rate_limter_strategy:
                status_code = await rate_limter_strategy.configuration(config=algorithms[algo])
                return status_code
        except Exception as e:
            return JSONResponse(
                content={
                    "responseData": {},
                    "message": [{"msg": f"Something went wrong {e}"}],
                    "success": False,
                    "code": HttpStatus.HTTP_400_BAD_REQUEST,
                },
                status_code=HttpStatus.HTTP_400_BAD_REQUEST,
            )


class CallerIdentityMiddlware(BaseHTTPMiddleware):
    async def get_caller_identity(request: Request):
        try:
            client_ip = request.client.host
            endpoint = request.url.path
            user = MockDB.get_db()
            user.update({
                'client_ip': client_ip,
                'endpoint': endpoint
            })
            kwargs = {"user_state": user}
            update_global_context(**kwargs)
        except Exception as e:
            return JSONResponse(
                content={
                    "responseData": {},
                    "message": [{"msg": f"Something went wrong {e}"}],
                    "success": False,
                    "code": HttpStatus.HTTP_400_BAD_REQUEST,
                },
                status_code=HttpStatus.HTTP_400_BAD_REQUEST,
            )
