from datetime import datetime
from typing import List

from fastapi import APIRouter, Depends

from api.dependency import filter_posts_parameters
from api.exceptions import HttpNotFoundException, HttpConflictException
from db.postgres import foreign_key_violation_exception
from models.post import Post
from models.thread import Thread, ThreadUpdate
from models.vote import Vote
from services.thread import thread_service

router = APIRouter()


@router.post('/{slug_or_id}/create', status_code=201)
async def create_post(
        slug_or_id: str,
        item: List[Post]
) -> List[Post]:
    thread = await thread_service.get_by_slug_or_id(slug_or_id)
    if not thread:
        
        raise HttpNotFoundException()

    if len(item) == 0:
        return item

    now = datetime.now()
    for post in item:
        post.created = now
        post.thread = thread.id
        post.forum = thread.forum

    posts, invalid_parents = await thread_service.create_posts(item)

    
    if invalid_parents:
        raise HttpConflictException()
    if not posts:
        raise HttpNotFoundException()

    return posts


@router.get('/{slug_or_id}/details')
async def get_thread_details(
        slug_or_id: str
) -> Thread:
    thread = await thread_service.get_by_slug_or_id(slug_or_id)

    
    if not thread:
        raise HttpNotFoundException()

    return thread


@router.post('/{slug_or_id}/details')
async def edit_thread(
        slug_or_id: str,
        item: ThreadUpdate
) -> Thread:
    thread = await thread_service.update_by_slug_or_id(slug_or_id, item)
    
    if not thread:
        raise HttpNotFoundException()

    return thread


@router.get('/{slug_or_id}/posts')
async def get_thread_posts(
        slug_or_id: str,
        filter_params: dict = Depends(filter_posts_parameters),
        sort: str = 'flat'
) -> List[Post]:
    thread = None
    if slug_or_id.isdigit():
        thread = await thread_service.get_by_id(int(slug_or_id))

    if not thread:
        thread = await thread_service.get_by_slug(slug_or_id)
    if not thread:
        
        raise HttpNotFoundException()

    posts = await thread_service.get_posts(thread.id, filter_params, sort)

    
    return posts



@router.post('/{slug_or_id}/vote')
async def vote_for_thread(
        slug_or_id: str,
        item: Vote
) -> Thread:
    if slug_or_id.isdigit():
        try:
            await thread_service.vote_by_id(int(slug_or_id), item)
        except foreign_key_violation_exception:
            pass
        else:
            thread = await thread_service.get_by_id(int(slug_or_id))
            return thread

    thread = await thread_service.get_by_slug(slug_or_id)
    if not thread:
        raise HttpNotFoundException()
    try:
        await thread_service.vote_by_id(thread.id, item)
    except foreign_key_violation_exception:
        raise HttpNotFoundException()
    thread = await thread_service.get_by_id(thread.id)

    if not thread:
        raise HttpNotFoundException()

    return thread
