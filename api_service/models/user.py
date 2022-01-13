import orjson
from pydantic import BaseModel

from models.common import orjson_dumps


class User(BaseModel):
    email: str
    fullname: str
    nickname: str = ''
    about: str

    class Config:
        json_loads = orjson.loads
        json_dumps = orjson_dumps


class UserUpdate(BaseModel):
    email: str = ''
    fullname: str = ''
    about: str = ''

    class Config:
        json_loads = orjson.loads
        json_dumps = orjson_dumps
