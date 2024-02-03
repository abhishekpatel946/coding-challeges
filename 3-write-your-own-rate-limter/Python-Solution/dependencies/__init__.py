import os

import yaml
from core import RateLimiterStrategyFactory


class RateLimiter:
    async def __init__(self):
        cwd = os.getcwd()
        with open(f"{cwd}/rate_limit.config.yml", "r") as f:
            self.config = yaml.safe_load(f)

        await RateLimiter.set_rate_limiter_configuration(self.config)

    async def set_rate_limiter_configuration(configuration):
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
