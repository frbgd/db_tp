from pydantic import BaseModel


class User(BaseModel):
    email: str
    fullname: str
    nickname: str = ''
    about: str


class UserUpdate(BaseModel):
    email: str
    fullname: str
    about: str
