import json
from config import settings
import redis


class RedisClient:
    def __init__(self):
        try:
            self.redis = redis.StrictRedis(host=settings.REDIS_ENDPOINT, port=settings.REDIS_PORT, db=0)
        except Exception as e:
            print(e)

    def get_client(self):
        return self.redis

    async def close(self):
        await self.redis.close()

    def set_json(self, key, data, expiry):
        value = None
        try:
            value = json.dumps(data)
            with self.redis.pipeline(transaction=True) as pipe:
                pipe.set(
                    key, value, ex=expiry
                ).execute()
        except Exception as e:
            print(f"Failed to set json for key {key}", e)

    def get_json(self, key):
        try:
            data = self.redis.get(key)
            return json.loads(data) if data else None
        except Exception as e:
            print(f"Error while reading value of key {key}", e)
