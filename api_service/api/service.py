from fastapi import APIRouter, Depends
from starlette.responses import Response

from models.status import Status
from services.service import ServiceService, get_service_service

router = APIRouter()


@router.post('/clear')
async def clear_database(
        service_service: ServiceService = Depends(get_service_service),
):
    await service_service.clear_database()

    await service_service.db.close()  # FIXME говнокод

    return Response()


@router.get('/status')
async def get_database_status(
        service_service: ServiceService = Depends(get_service_service)
) -> Status:
    status = await service_service.get_status()
    await service_service.db.close()  # FIXME говнокод
    return status
