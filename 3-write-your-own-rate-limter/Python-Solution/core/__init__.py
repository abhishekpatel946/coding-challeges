from abc import ABCMeta, abstractmethod
from redis import RedisClient
from context import _global_state, update_global_context, delete_global_context


class RateLimitStrategy(metaclass=ABCMeta):
    @abstractmethod
    def authenticate(self, **kwargs):
        raise NotImplementedError()


class TokenBucketStrategy:
    async def configuration(self, config={}):
        try:
            kwargs = {"token_bucket_configuration": config}
            update_global_context(**kwargs)
            print("Configuring token bucket configuration", config, _global_state)
            delete_global_context(key="token_bucket_configuration")
            print("reseting token bucket configuration", config, _global_state)
            # redis_client = RedisClient()
            # await redis_client.set_json(key='capacity',
            #                             data=config.get('capacity', 0), expiry=config.get('refresh_interval', 0))
        except Exception as e:
            print(e)


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
