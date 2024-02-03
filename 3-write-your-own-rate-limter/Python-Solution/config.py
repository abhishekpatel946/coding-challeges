from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    APP_NAME: str = 'FastAPI application'
    ENVIRONMENT: str
    REDIS_ENDPOINT: str
    REDIS_USERNAME: str

    class Config:
        env_file = ".env"


settings = Settings()
