from http import HTTPStatus

from fastapi import APIRouter, Depends
from starlette.responses import JSONResponse

from api.exceptions import HttpNotFoundException, HttpConflictException
from models.user import User, UserUpdate
from services.user import UserService, get_user_service

router = APIRouter()


@router.post('/{nickname}/create', status_code=201)
async def create_user(
        nickname: str,
        item: User,
        user_service: UserService = Depends(get_user_service)
):
    # TODO валидация
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
        await user_service.db.close()  # FIXME говнокод
        return response

    await user_service.db.close()  # FIXME говнокод
    return users[0]


@router.get('/{nickname}/profile')
async def get_user_details(
        nickname: str,
        user_service: UserService = Depends(get_user_service)
) -> User:
    user = await user_service.get_by_nickname(nickname)

    await user_service.db.close()  # FIXME говнокод
    if not user:
        raise HttpNotFoundException()

    return user


@router.post('/{nickname}/profile')
async def edit_user_details(
        nickname: str,
        item: UserUpdate,
        user_service: UserService = Depends(get_user_service)
) -> User:
    # TODO валидация
    user, not_unique = await user_service.update_by_nickname(nickname, item)

    await user_service.db.close()  # FIXME говнокод
    if not_unique:
        raise HttpConflictException()
    if not user:
        raise HttpNotFoundException()

    return user
