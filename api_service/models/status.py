from pydantic import BaseModel


class Status(BaseModel):
    forum: int
    post: int
    thread: int
    user: int
