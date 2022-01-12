from pydantic import BaseModel


class Forum(BaseModel):
    slug: str
    title: str
    user: str
    posts: int = 0
    threads: int = 0
