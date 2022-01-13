from typing import Optional, Tuple, List

from db.postgres import PostgressDbEngine, item_not_unique_exception
from models.user import User, UserUpdate


class UserService:
    def __init__(self, db: PostgressDbEngine):
        self.db = db

    async def get_users_by_nickname_or_email(self, nickname: str, email: str) -> List[User]:
        users = []
        user_by_nickname = await self.get_by_nickname(nickname)
        if user_by_nickname:
            users.append(user_by_nickname)
        user_by_email = await self.get_by_email(email)
        if user_by_email:
            if not user_by_nickname or user_by_nickname and user_by_email.nickname != user_by_nickname.nickname:
                users.append(user_by_email)
        return users

    async def create_user(self, item: User) -> Tuple[List[User], bool]:
        try:
            await self.db.insert(
                """INSERT INTO users (nickname, email, fullname, about)
				VALUES ($1, $2, $3, $4)""",
                item.nickname, item.email, item.fullname, item.about
            )
        except item_not_unique_exception:
            return await self.get_users_by_nickname_or_email(item.nickname, item.email), True

        return [item], False

    async def get_by_email(self, email: str) -> Optional[User]:
        value = await self.db.get_one(
            """SELECT 	nickname,
                        email,
                        fullname,
                        about
                FROM users
                WHERE email = $1""",
            email
        )

        if not value:
            return

        return User(
            email=value['email'],
            nickname=value['nickname'],
            fullname=value['fullname'],
            about=value['about']
        )

    async def get_by_nickname(self, nickname: str) -> Optional[User]:
        value = await self.db.get_one(
            """SELECT 	nickname,
                        email,
                        fullname,
                        about
                FROM users
                WHERE nickname = $1""",
            nickname
        )

        if not value:
            return

        return User(
            email=value['email'],
            nickname=value['nickname'],
            fullname=value['fullname'],
            about=value['about']
        )

    async def update_by_nickname(self, nickname: str, item: UserUpdate) -> Tuple[Optional[User], bool]:
        if item.email:
            user_email = item.email
        else:
            user_email = None
        if item.fullname:
            user_fullname = item.fullname
        else:
            user_fullname = None
        if item.about:
            user_about = item.about
        else:
            user_about = None
        try:
            value = await self.db.update(
                """UPDATE users
                    SET email=COALESCE(NULLIF($2, ''), email),
                        fullname=COALESCE(NULLIF($3, ''), fullname),
                        about=COALESCE(NULLIF($4, ''), about)
                    WHERE nickname = $1
                    RETURNING nickname, email, fullname, about
                """,
                nickname, user_email, user_fullname, user_about)
        except item_not_unique_exception:
            return None, True

        if not value:
            return None, False

        return User(
            email=value['email'],
            fullname=value['fullname'],
            nickname=value['nickname'],
            about=value['about']
        ), False


user_service: UserService = None


def get_user_service() -> UserService:
    return user_service
