import json
from config import settings
import aioredis


class RedisClient:
    def __init__(self):
        try:
            self.redis = aioredis.from_url(f"redis://{settings.REDIS_ENDPOINT}")
        except Exception as e:
            print(e)

    def get_client(self):
        return self.redis

    async def close(self):
        await self.redis.close()

    async def set_json(self, key, data, expiry):
        value = None
        try:
            value = json.dumps(data)
            async with self.redis.pipeline(transaction=True) as pipe:
                await pipe.set(
                    key, value, ex=expiry
                ).execute()
        except Exception as e:
            print(f"Failed to set json for key {key}", e)

    async def get_json(self, key):
        try:
            data = await self.redis.get(key)
            return json.loads(data) if data else None
        except Exception as e:
            print(f"Error while reading value of key {key}", e)
