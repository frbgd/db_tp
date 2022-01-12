from datetime import datetime
from typing import Optional

from pydantic import BaseModel

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


class FullPost(BaseModel):
    post: Post
    author: User = None
    thread: Thread = None
    forum: Forum = None


class PostUpdate(BaseModel):
    message: str = ''
