import orjson
from pydantic import BaseModel

from models.common import orjson_dumps


class Vote(BaseModel):
    nickname: str
    voice: int

    class Config:
        json_loads = orjson.loads
        json_dumps = orjson_dumps
