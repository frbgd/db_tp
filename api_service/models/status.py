import orjson
from pydantic import BaseModel

from models.common import orjson_dumps


class Status(BaseModel):
    forum: int
    post: int
    thread: int
    user: int

    class Config:
        json_loads = orjson.loads
        json_dumps = orjson_dumps
