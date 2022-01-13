from http import HTTPStatus
from typing import List

from fastapi import APIRouter, Depends, Response

from api.dependency import filter_parameters, filter_users_parameters
from api.exceptions import HttpNotFoundException
from models.forum import Forum
from models.thread import Thread
from models.user import User
from services.forum import forum_service

router = APIRouter()


@router.post('/create', status_code=201)
async def create_forum(
        item: Forum,
        response: Response
):
    forum, not_unique = await forum_service.create_forum(item)

    if not forum:
        raise HttpNotFoundException()
    if not_unique:
        response.status_code = HTTPStatus.CONFLICT

    return forum


@router.get('/{slug}/details')
async def get_forum_details(
        slug: str
) -> Forum:
    forum = await forum_service.get_by_slug(slug)

    if not forum:
        raise HttpNotFoundException()

    return forum


@router.post('/{slug}/create', status_code=201)
async def create_thread(
        slug: str,
        item: Thread,
        response: Response
):
    if not item.forum:
        item.forum = slug
    thread, not_unique = await forum_service.create_thread(item)

    if not thread:
        raise HttpNotFoundException()
    if not_unique:
        response.status_code = HTTPStatus.CONFLICT

    return thread


@router.get('/{slug}/users')
async def get_forum_users(
        slug: str,
        filter_params: dict = Depends(filter_users_parameters)
) -> List[User]:
    users = await forum_service.get_forum_users(
        slug=slug,
        filter_params=filter_params
    )

    if not isinstance(users, list):
        raise HttpNotFoundException()

    return users


@router.get('/{slug}/threads')
async def get_forum_threads(
        slug: str,
        filter_params: dict = Depends(filter_parameters)
) -> List[Thread]:
    threads = await forum_service.get_forum_threads(
        slug=slug,
        filter_params=filter_params
    )

    if not isinstance(threads, list):
        raise HttpNotFoundException()

    return threads
