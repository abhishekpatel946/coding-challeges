from abc import ABCMeta, abstractmethod

from fastapi import status as HttpStatus
from starlette.responses import JSONResponse

from context import _global_state, update_global_context
from utils.redis import RedisClient


class RateLimitStrategy(metaclass=ABCMeta):
    @abstractmethod
    def authenticate(self, **kwargs):
        raise NotImplementedError()


class TokenBucketStrategy:
    async def configuration(self, config={}):
        try:
            kwargs = {"token_bucket_configuration": config}
            update_global_context(**kwargs)
            user_id = _global_state.get('user_state', {}).get('user_id', None)
            endpoint = config.get('endpoints', {}).get('endpoint', None)
            if endpoint and user_id:
                redis_client = RedisClient()
                cache = redis_client.get_json(key=f'{user_id}#{endpoint}')
                if not cache:
                    # set the configuration in the cache
                    redis_client.set_json(key=f'{user_id}#{endpoint}',
                                          data={
                                              'endpoint': endpoint,
                                              'capacity': config.get('endpoints', {}).get('limit', 0),
                                          }, expiry=config.get('endpoints', {}).get('interval', 0))
                    _global_state.get('user_state', {}).update({
                        'capacity': config.get('endpoints', {}).get('limit', 0)
                    })
                    kwargs = {"user_state": _global_state.get(
                        'user_state', {})}
                    update_global_context(**kwargs)
                else:
                    # update the configuration
                    limit = cache.get('capacity', 0)
                    if limit:
                        limit -= 1
                        cache.update({'capacity': limit})
                        redis_client.set_json(key=f'{user_id}#{endpoint}', data=cache, expiry=config.get(
                            'endpoints', {}).get('interval', 0))
                        _global_state.get('user_state', {}).update(
                            {'capacity': limit})
                        kwargs = {"user_state": _global_state.get(
                            'user_state', {})}
                        update_global_context(**kwargs)
                        return HttpStatus.HTTP_202_ACCEPTED
                    else:
                        return HttpStatus.HTTP_429_TOO_MANY_REQUESTS
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


class LeakyBucketMirrorStrategy:
    pass


class LeakyBucketQueueStrategy:
    pass


class SlidingBucketLogStrategy:
    pass


class SlidingBucketCounterStrategy:
    pass


class RateLimiterStrategyFactory(RateLimitStrategy):
    @staticmethod
    def get_strategy(approach):
        if approach == 'token_bucket':
            return TokenBucketStrategy()
        elif approach == 'leaky_bucket_mirror':
            return LeakyBucketMirrorStrategy()
        elif approach == 'leaky_bucket_queue':
            return LeakyBucketQueueStrategy()
        elif approach == 'sliding_bucket_log':
            return SlidingBucketLogStrategy()
        elif approach == 'sliding_bucket_counter':
            return SlidingBucketCounterStrategy()
        else:
            raise ValueError("Invalid algorithm or strategy")
