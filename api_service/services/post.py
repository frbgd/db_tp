from typing import Optional

from db.postgres import PostgressDbEngine
from models.forum import Forum
from models.post import Post, FullPost, PostUpdate
from models.thread import Thread
from models.user import User


class PostService:
    def __init__(self, db: PostgressDbEngine):
        self.db = db

    async def get_user_by_nickname(self, nickname: str) -> Optional[User]:
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

    async def get_thread_by_id(self, id_: int) -> Optional[Thread]:
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

    async def get_forum_by_slug(self, slug: str) -> Optional[Forum]:
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

    async def get_by_id(self, id_: int) -> Optional[FullPost]:
        value = await self.db.get_one(
            """SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE id = $1""",
            id_
        )

        if not value:
            return

        # TODO переместить в соответствующие сервисы
        author = await self.get_user_by_nickname(value['author_nickname'])
        thread = await self.get_thread_by_id(value['thread_id'])
        forum = await self.get_forum_by_slug(value['forum_slug'])

        return FullPost(
            post=Post(
                id=value['id'],
                message=value['message'],
                created=value['created'],
                author=value['author_nickname'],
                forum=value['forum_slug'],
                isEdited=value['is_edited'],
                parent=value['parent'],
                thread=value['thread_id']
            ),
            author=author,
            thread=thread,
            forum=forum
        )

    async def update_by_id(self, id_: int, item: PostUpdate) -> Optional[Post]:
        post_message = item.message if item.message else None

        value = await self.db.update(
            """UPDATE posts
			SET message   = COALESCE($2, message),
				is_edited = CASE
								WHEN (is_edited = TRUE
									OR (is_edited = FALSE AND $2 IS NOT NULL AND $2 <> message)) THEN TRUE
								ELSE FALSE
					END
			WHERE id = $1
			RETURNING id, thread_id, author_nickname, forum_slug, is_edited, message, parent, created
			""",
        id_, post_message)

        if not value:
            return

        return Post(
            id=value['id'],
            message=value['message'],
            created=value['created'],
            author=value['author_nickname'],
            forum=value['forum_slug'],
            isEdited=value['is_edited'],
            parent=value['parent'],
            thread=value['thread_id']
        )


post_service: PostService = None


def get_post_service() -> PostService:
    return post_service
