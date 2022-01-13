from datetime import datetime
from typing import Optional

import orjson
from pydantic import BaseModel

from models.common import orjson_dumps


class Thread(BaseModel):
    id: int = 0
    slug: Optional[str] = None
    title: str
    message: str
    created: datetime = datetime.now()
    author: str
    forum: str = ''
    votes: int = 0

    class Config:
        json_loads = orjson.loads
        json_dumps = orjson_dumps


class ThreadUpdate(BaseModel):
    title: str = ''
    message: str = ''

    class Config:
        json_loads = orjson.loads
        json_dumps = orjson_dumps
