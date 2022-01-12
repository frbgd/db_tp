from pydantic import BaseModel


class Vote(BaseModel):
    nickname: str
    voice: int
