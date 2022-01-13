from typing import Optional, List, Tuple

from db.postgres import PostgressDbEngine, item_not_unique_exception, not_null_constraint_exception
from models.forum import Forum
from models.thread import Thread
from models.user import User

get_users_by_forum_slug_since_sql = {
    True: """SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		  AND user_nickname < $2
		ORDER BY user_nickname DESC
		LIMIT $3""",
    False: """SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		  AND user_nickname > $2
		ORDER BY user_nickname
		LIMIT $3"""
}

get_users_by_forum_slug_sql = {
    True: """SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		ORDER BY user_nickname DESC
		LIMIT $2""",
    False: """SELECT nickname,
			   email,
			   fullname,
			   about
		FROM forums_users
		WHERE forum_slug = $1
		ORDER BY user_nickname
		LIMIT $2"""
}

get_threads_by_forum_slug_since_sql = {
    True: """SELECT id,
			   slug,
			   forum_slug,
			   user_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1 AND created <= $2
		ORDER BY created DESC
		LIMIT $3""",
    False: """SELECT id,
			   slug,
			   forum_slug,
			   user_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1 AND created >= $2
		ORDER BY created
		LIMIT $3"""
}

get_threads_by_forum_slug_sql = {
    True: """SELECT id,
			   slug,
			   forum_slug,
			   user_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1
		ORDER BY created DESC
		LIMIT $2""",
    False: """SELECT id,
			   slug,
			   forum_slug,
			   user_nickname,
			   title,
			   message,
			   votes,
			   created
		FROM threads
		WHERE forum_slug = $1
		ORDER BY created
		LIMIT $2"""
}


class ForumService:
    def __init__(self, db: PostgressDbEngine):
        self.db = db

    async def get_by_slug(self, slug: str) -> Optional[Forum]:
        value = await self.db.get_one(
            """SELECT slug, title, threads, posts, owner_nickname
			    FROM forums
				WHERE slug = $1""",
            slug
        )

        if not value:
            return

        return Forum(
            slug=value['slug'],
            title=value['title'],
            user=value['owner_nickname'],
            posts=value['posts'],
            threads=value['threads']
        )

    async def get_thread_by_slug(self, slug: str) -> Optional[Thread]:
        # TODO переместить в другой сервис
        value = await self.db.get_one(
            """SELECT id,
				   slug,
				   forum_slug,
				   user_nickname,
				   title,
				   message,
				   votes,
				   created
			FROM threads
			WHERE slug = $1""",
            slug
        )

        if not value:
            return

        return Thread(
            id=value['id'],
            slug=value['slug'],
            title=value['title'],
            message=value['message'],
            created=value['created'],
            author=value['user_nickname'],
            forum=value['forum_slug'],
            votes=value['votes']
        )

    async def get_forum_users(self, slug, filter_params: dict) -> Optional[List[User]]:
        since = filter_params.get('since')
        limit = filter_params.get('limit')
        desc = filter_params.get('desc')

        if since:
            sql = get_users_by_forum_slug_since_sql[desc]
            values = await self.db.get_many(
                sql,
                slug, since, limit
            )
        else:
            sql = get_users_by_forum_slug_sql[desc]
            values = await self.db.get_many(
                sql,
                slug, limit
            )

        if len(values) == 0:
            return [] if await self.get_by_slug(slug) else None
        else:
            return [
                User(
                    email=value['email'],
                    fullname=value['fullname'],
                    nickname=value['nickname'],
                    about=value['about']
                ) for value in values
            ]

    async def get_forum_threads(self, slug: str, filter_params: dict) -> Optional[List[Thread]]:
        since = filter_params.get('since')
        limit = filter_params.get('limit')
        desc = filter_params.get('desc')

        if since:
            sql = get_threads_by_forum_slug_since_sql[desc]
            values = await self.db.get_many(
                sql,
                slug, since, limit
            )
        else:
            sql = get_threads_by_forum_slug_sql[desc]
            values = await self.db.get_many(
                sql,
                slug, limit
            )

        if len(values) == 0:
            return [] if await self.get_by_slug(slug) else None
        else:
            return [
                Thread(
                    id=value['id'],
                    slug=value['slug'],
                    title=value['title'],
                    message=value['message'],
                    created=value['created'],
                    author=value['user_nickname'],
                    forum=value['forum_slug'],
                    votes=value['votes']
                ) for value in values
            ]

    async def create_forum(self, item: Forum) -> Tuple[Optional[Forum], bool]:
        try:
            value = await self.db.insert(
                """INSERT INTO forums (slug, title, threads, posts, owner_nickname)
                    VALUES ($1, $2, $3, $4, (
                                SELECT nickname FROM users WHERE nickname = $5
                    )) RETURNING owner_nickname""",
                item.slug, item.title, item.threads, item.posts, item.user
            )
        except item_not_unique_exception:
            return await self.get_by_slug(item.slug), True
        except not_null_constraint_exception:
            return None, False

        item.user = value['owner_nickname']
        return item, False

    async def create_thread(self, item: Thread) -> Tuple[Optional[Thread], bool]:
        try:
            value = await self.db.insert(
                """INSERT INTO threads (slug, user_nickname, title, message, created, forum_slug)
			VALUES ($1, (SELECT nickname FROM users WHERE nickname = $2), $3, $4, $5,
					(SELECT slug FROM forums WHERE slug = $6))
			RETURNING id, user_nickname, forum_slug, created""",
                item.slug, item.author, item.title, item.message, item.created, item.forum
            )
        except item_not_unique_exception:
            return await self.get_thread_by_slug(item.slug), True
        except not_null_constraint_exception:
            return None, False

        item.id = value['id']
        item.author = value['user_nickname']
        item.forum = value['forum_slug']
        return item, False


forum_service: ForumService = None


def get_forum_service() -> ForumService:
    return forum_service
