from datetime import datetime
from typing import Optional

import orjson
from pydantic import BaseModel

from models.common import orjson_dumps
from models.forum import Forum
from models.thread import Thread
from models.user import User


class Post(BaseModel):
    id: int = 0
    message: str
    created: datetime = None
    author: str
    forum: str = ''
    isEdited: bool = False
    parent: Optional[int] = None
    thread: int = 0

    class Config:
        # Заменяем стандартную работу с json на более быструю
        json_loads = orjson.loads
        json_dumps = orjson_dumps


class FullPost(BaseModel):
    post: Post
    author: User = None
    thread: Thread = None
    forum: Forum = None

    class Config:
        # Заменяем стандартную работу с json на более быструю
        json_loads = orjson.loads
        json_dumps = orjson_dumps


class PostUpdate(BaseModel):
    message: str = ''

    class Config:
        # Заменяем стандартную работу с json на более быструю
        json_loads = orjson.loads
        json_dumps = orjson_dumps
