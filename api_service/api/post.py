from fastapi import APIRouter, Depends

from api.exceptions import HttpNotFoundException
from models.post import Post, FullPost, PostUpdate
from services.post import PostService, get_post_service

router = APIRouter()


@router.get('/{id}/details')
async def get_post_details(
        id: int,
        post_service: PostService = Depends(get_post_service)
) -> FullPost:
    # TODO related поле
    post = await post_service.get_by_id(id)

    await post_service.db.close()  # FIXME говнокод

    if not post:
        raise HttpNotFoundException()

    return post


@router.post('/{id}/details')
async def edit_post(
        id: int,
        item: PostUpdate,
        post_service: PostService = Depends(get_post_service)
):
    # TODO валидация
    post = await post_service.update_by_id(id, item)

    await post_service.db.close()  # FIXME говнокод

    if not post:
        raise HttpNotFoundException()

    return post
