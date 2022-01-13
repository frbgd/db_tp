from http import HTTPStatus

from fastapi import APIRouter
from starlette.responses import JSONResponse

from api.exceptions import HttpNotFoundException, HttpConflictException
from models.user import User, UserUpdate
from services.user import user_service

router = APIRouter()


@router.post('/{nickname}/create', status_code=201)
async def create_user(
        nickname: str,
        item: User,
):
    item.nickname = nickname
    users, not_unique = await user_service.create_user(item)
    if not_unique:
        # FIXME очень говнокод, при наличии времени сделать нормально
        response = JSONResponse(content=[
            {
                'email': user.email,
                'fullname': user.fullname,
                'nickname': user.nickname,
                'about': user.about
            } for user in users
        ])
        response.status_code = HTTPStatus.CONFLICT
        return response
    return users[0]


@router.get('/{nickname}/profile')
async def get_user_details(
        nickname: str,
) -> User:
    user = await user_service.get_by_nickname(nickname)
    if not user:
        raise HttpNotFoundException()

    return user


@router.post('/{nickname}/profile')
async def edit_user_details(
        nickname: str,
        item: UserUpdate,
) -> User:
    user, not_unique = await user_service.update_by_nickname(nickname, item)
    if not_unique:
        raise HttpConflictException()
    if not user:
        raise HttpNotFoundException()

    return user
