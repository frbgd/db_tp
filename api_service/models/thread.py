from datetime import datetime
from typing import Optional

from pydantic import BaseModel


class Thread(BaseModel):
    id: int = 0
    slug: Optional[str] = None
    title: str
    message: str
    created: datetime = datetime.now()
    author: str
    forum: str = ''
    votes: int = 0


class ThreadUpdate(BaseModel):
    title: str
    message: str
