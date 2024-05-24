from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    APP_NAME: str = "FastAPI distributed UID generator"
    ENVIRONMENT: str
    REDIS_ENDPOINT: str
    REDIS_USERNAME: str
    REDIS_PORT: int

    class Config:
        env_file = ".env"


settings = Settings()
