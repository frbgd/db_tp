from fastapi import Depends

from db.postgres import PostgressDbEngine, get_postgres
from models.status import Status


class ServiceService:
    def __init__(self, db: PostgressDbEngine):
        self.db = db

    async def get_status(self) -> Status:
        user_value = await self.db.get_one('SELECT COUNT(*) FROM users;')
        forum_value = await self.db.get_one('SELECT COUNT(*) FROM forums;')
        thread_value = await self.db.get_one('SELECT COUNT(*) FROM threads;')
        post_value = await self.db.get_one('SELECT COUNT(*) FROM posts;')

        return Status(
            user=user_value['count'],
            forum=forum_value['count'],
            thread=thread_value['count'],
            post=post_value['count']
        )

    async def clear_database(self):
        await self.db.execute('TRUNCATE users, forums, threads, votes, posts, forums_users_nicknames')


service_service: ServiceService = None


def get_service_service() -> ServiceService:
    return service_service
