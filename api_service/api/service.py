from fastapi import APIRouter, Depends
from starlette.responses import Response

from models.status import Status
from services.service import service_service

router = APIRouter()


@router.post('/clear')
async def clear_database(
):
    await service_service.clear_database()

    await service_service.db.close()  # FIXME говнокод

    return Response()


@router.get('/status')
async def get_database_status(
) -> Status:
    status = await service_service.get_status()
    await service_service.db.close()  # FIXME говнокод
    return status
