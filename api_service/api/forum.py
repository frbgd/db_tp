from http import HTTPStatus
from typing import List

from fastapi import APIRouter, Depends, Response

from api.dependency import filter_parameters, filter_users_parameters
from api.exceptions import HttpNotFoundException
from models.forum import Forum
from models.thread import Thread
from models.user import User
from services.forum import ForumService, get_forum_service

router = APIRouter()


@router.post('/create', status_code=201)
async def create_forum(
        item: Forum,
        response: Response,
        forum_service: ForumService = Depends(get_forum_service)
):
    # TODO валидировать входные данные форума (числа, slug)
    forum, not_unique = await forum_service.create_forum(item)

    await forum_service.db.close()  # FIXME говнокод

    if not forum:
        raise HttpNotFoundException()
    if not_unique:
        response.status_code = HTTPStatus.CONFLICT

    return forum


@router.get('/{slug}/details')
async def get_forum_details(
        slug: str,
        forum_service: ForumService = Depends(get_forum_service)
) -> Forum:
    forum = await forum_service.get_by_slug(slug)

    await forum_service.db.close()  # FIXME говнокод

    if not forum:
        raise HttpNotFoundException()

    return forum


@router.post('/{slug}/create', status_code=201)
async def create_thread(
        slug: str,
        item: Thread,
        response: Response,
        forum_service: ForumService = Depends(get_forum_service)
):
    # TODO валидировать входные данные топика (slug, created)
    # TODO возможно надо title переводить в slug
    if not item.forum:
        item.forum = slug
    thread, not_unique = await forum_service.create_thread(item)

    await forum_service.db.close()  # FIXME говнокод

    if not thread:
        raise HttpNotFoundException()
    if not_unique:
        response.status_code = HTTPStatus.CONFLICT

    return thread


# TODO оттестить параметры
@router.get('/{slug}/users')
async def get_forum_users(
        slug: str,
        filter_params: dict = Depends(filter_users_parameters),
        forum_service: ForumService = Depends(get_forum_service)
) -> List[User]:
    users = await forum_service.get_forum_users(
        slug=slug,
        filter_params=filter_params
    )

    await forum_service.db.close()  # FIXME говнокод

    if not isinstance(users, list):
        raise HttpNotFoundException()

    return users


# TODO оттестить параметры
@router.get('/{slug}/threads')
async def get_forum_threads(
        slug: str,
        filter_params: dict = Depends(filter_parameters),
        forum_service: ForumService = Depends(get_forum_service)
) -> List[Thread]:
    threads = await forum_service.get_forum_threads(
        slug=slug,
        filter_params=filter_params
    )

    await forum_service.db.close()  # FIXME говнокод

    if not isinstance(threads, list):
        raise HttpNotFoundException()

    return threads
