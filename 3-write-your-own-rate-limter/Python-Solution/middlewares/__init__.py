import os

import yaml
from fastapi import Request
from fastapi import status as HttpStatus
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.responses import JSONResponse

from core import RateLimiterStrategyFactory


class RateLimiterMiddleware(BaseHTTPMiddleware):
    def __init__(self, app):
        super().__init__(app)

    async def dispatch(self, request: Request, call_next):
        try:
            cwd = os.getcwd()
            with open(f"{cwd}/rate_limit.config.yml", "r") as f:
                self.config = yaml.safe_load(f)

            return await RateLimiterMiddleware.set_rate_limiter_configuration(self.config)

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

    async def set_rate_limiter_configuration(configuration):
        try:
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
                await rate_limter_strategy.configuration(config=algorithms[algo])
        except Exception as e:
            print(e)
