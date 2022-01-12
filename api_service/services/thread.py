from typing import Optional, List, Tuple

from asyncpg import RaiseError
from fastapi import Depends

from db.postgres import PostgressDbEngine, get_postgres, not_null_constraint_exception, foreign_key_violation_exception
from models.post import Post
from models.thread import Thread, ThreadUpdate
from models.vote import Vote


get_posts_since_sql = {
    True: {
        'flat': """SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND id < $2
			ORDER BY id DESC
			LIMIT $3""",
        'tree': """SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path < (SELECT path FROM posts WHERE id = $2)
			ORDER BY path DESC
			LIMIT $3""",
        'parent_tree': """WITH roots AS (
				SELECT DISTINCT path[1]
				FROM posts
				WHERE thread_id = $1
				  AND parent IS NULL
				  AND path[1] < (SELECT path[1] FROM posts WHERE id = $2)
         		ORDER BY path[1] DESC
				LIMIT $3
			)
			SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path[1] IN (SELECT * FROM roots)
			ORDER BY path[1] DESC, path[2:]"""
    },
    False: {
        'flat': """SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND id > $2
			ORDER BY id
			LIMIT $3""",
        'tree': """SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path > (SELECT path FROM posts WHERE id = $2)
			ORDER BY path
			LIMIT $3""",
        'parent_tree': """WITH roots AS (
				SELECT DISTINCT path[1]
				FROM posts
				WHERE thread_id = $1
				  AND parent IS NULL
				  AND path[1] > (SELECT path[1] FROM posts WHERE id = $2)
         		ORDER BY path[1]
				LIMIT $3
			)
			SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path[1] IN (SELECT * FROM roots)
			ORDER BY path"""
    }
}

get_posts_sql = {
    True: {
        'flat': """SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			ORDER BY id DESC
			LIMIT $2""",
        'tree': """SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			ORDER BY path DESC
			LIMIT $2""",
        'parent_tree': """WITH roots AS (
				SELECT DISTINCT path[1]
				FROM posts
				WHERE thread_id = $1
				ORDER BY path[1] DESC
				LIMIT $2
			)
			SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path[1] IN (SELECT * FROM roots)
			ORDER BY path[1] DESC, path[2:]"""
    },
    False: {
        'flat': """SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			ORDER BY id
			LIMIT $2""",
        'tree': """SELECT id,
					   thread_id,
					   author_nickname,
					   forum_slug,
					   is_edited,
					   message,
					   parent,
					   created
				FROM posts
				WHERE thread_id = $1
				ORDER BY path
				LIMIT $2""",
        'parent_tree': """WITH roots AS (
				SELECT DISTINCT path[1]
				FROM posts
				WHERE thread_id = $1
				ORDER BY path[1]
				LIMIT $2
			)
			SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE thread_id = $1
			  AND path[1] IN (SELECT * FROM roots)
			ORDER BY path"""
    }
}


class ThreadService:
    def __init__(self, db: PostgressDbEngine):
        self.db = db

    async def get_by_slug_or_id(self, slug_or_id: str) -> Optional[Thread]:
        thread = None
        if slug_or_id.isdigit():
            thread = await self.get_by_id(slug_or_id)
        if thread:
            return thread
        return await self.get_by_slug(slug_or_id)

    async def get_by_slug(self, slug: str) -> Optional[Thread]:
        value = await self.db.get_one(
            """SELECT id,
				   slug,
				   forum_slug,
				   author_nickname,
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
            author=value['author_nickname'],
            forum=value['forum_slug'],
            votes=value['votes']
        )

    async def get_by_id(self, id_: str) -> Optional[Thread]:
        value = await self.db.get_one(
            """SELECT id,
                   slug,
                   forum_slug,
                   author_nickname,
                   title,
                   message,
                   votes,
                   created
            FROM threads
            WHERE id = $1""",
            id_
        )

        if not value:
            return

        return Thread(
            id=value['id'],
            slug=value['slug'],
            title=value['title'],
            message=value['message'],
            created=value['created'],
            author=value['author_nickname'],
            forum=value['forum_slug'],
            votes=value['votes']
        )

    async def create_posts(self, posts: List[Post]) -> Tuple[Optional[List[Post]], bool]:
        result_posts = []
        # TODO транзакционность
        for post in posts:
            try:
                value = await self.db.insert(
                    """INSERT INTO posts (thread_id, author_nickname, forum_slug, message, parent, created)
                    VALUES ($1, $2, $3, $4, $5, $6)
                    RETURNING id, thread_id, author_nickname, forum_slug, is_edited, message, parent, created""",
                    post.thread, post.author, post.forum, post.message, post.parent, post.created
                )
            except (not_null_constraint_exception, foreign_key_violation_exception):
                return None, False
            except RaiseError as exc:
                return None, True
            result_posts.append(
                Post(
                    id=value['id'],
                    thread=value['thread_id'],
                    author=value['author_nickname'],
                    forum=value['forum_slug'],
                    isEdited=value['is_edited'],
                    message=value['message'],
                    parent=value['parent'],
                    created=value['created']
                )
            )
        return result_posts, False

    async def update_by_slug_or_id(self, slug_or_id: str, item: ThreadUpdate) -> Optional[Thread]:
        thread = None
        if slug_or_id.isdigit():
            thread = await self.update_by_id(slug_or_id, item)
        if thread:
            return thread
        return await self.update_by_slug(slug_or_id, item)

    async def update_by_slug(self, slug: str, item: ThreadUpdate) -> Optional[Thread]:
        thread_title = item.title if item.title else None
        thread_message = item.message if item.message else None

        value = await self.db.update(
                """UPDATE threads
			SET title   = COALESCE(NULLIF($2, ''), title),
				message = COALESCE(NULLIF($3, ''), message)
			WHERE slug = $1
			RETURNING id, slug, forum_slug, author_nickname, title, message, votes, created""",
                slug, thread_title, thread_message
        )
        if not value:
            return

        return Thread(
            id=value['id'],
            slug=value['slug'],
            title=value['title'],
            message=value['message'],
            created=value['created'],
            author=value['author_nickname'],
            forum=value['forum_slug'],
            votes=value['votes']
        )

    async def update_by_id(self, id_: str, item: ThreadUpdate) -> Optional[Thread]:
        thread_title = item.title if item.title else None
        thread_message = item.message if item.message else None

        value = await self.db.update(
            """UPDATE threads
			SET title   = COALESCE(NULLIF($2, ''), title),
				message = COALESCE(NULLIF($3, ''), message)
			WHERE id = $1
			RETURNING id, slug, forum_slug, author_nickname, title, message, votes, created""",
            id_, thread_title, thread_message
        )
        if not value:
            return

        return Thread(
            id=value['id'],
            slug=value['slug'],
            title=value['title'],
            message=value['message'],
            created=value['created'],
            author=value['author_nickname'],
            forum=value['forum_slug'],
            votes=value['votes']
        )

    async def vote_by_id(self, id_: str, item: Vote):
        # FIXME я хз, тут только так работает
        nickname = item.nickname
        voice = item.voice
        # ¯\_(ツ)_/¯
        await self.db.execute(
            """INSERT INTO votes (thread_id, author_nickname, voice)
			VALUES ($1, $2, $3)
			ON CONFLICT (thread_id, author_nickname) DO UPDATE SET voice = $3""",
            id_, nickname, voice
        )

    async def get_posts(self, thread_id: int, filter_params: dict, sort: str) -> List[Post]:
        since = filter_params.get('since')
        limit = filter_params.get('limit')
        desc = filter_params.get('desc')

        if since:
            sql = get_posts_since_sql[desc][sort]
            values = await self.db.get_many(
                sql,
                thread_id, since, limit
            )
        else:
            sql = get_posts_sql[desc][sort]
            values = await self.db.get_many(
                sql,
                thread_id, limit
            )

        return [
            Post(
                id=value['id'],
                thread=value['thread_id'],
                author=value['author_nickname'],
                forum=value['forum_slug'],
                isEdited=value['is_edited'],
                message=value['message'],
                parent=value['parent'],
                created=value['created']
            ) for value in values
        ]


def get_thread_service(
        db: PostgressDbEngine = Depends(get_postgres)
) -> ThreadService:
    return ThreadService(db)
