import orjson
from pydantic import BaseModel

from models.common import orjson_dumps


class Forum(BaseModel):
    slug: str
    title: str
    user: str
    posts: int = 0
    threads: int = 0

    class Config:
        json_loads = orjson.loads
        json_dumps = orjson_dumps
