from datetime import datetime
from typing import Optional, Union

from fastapi import Query


async def filter_parameters(
    since: Optional[datetime] = Query(default=None),
    desc: Optional[bool] = Query(default=False),
    limit: Optional[int] = Query(default=100, ge=1),
):
    return {'limit': limit, 'since': since, 'desc': desc}


async def filter_posts_parameters(
        since: Optional[int] = Query(default=None),
        desc: Optional[bool] = Query(default=False),
        limit: Optional[int] = Query(default=100, ge=1),
):
    return {'limit': limit, 'since': since, 'desc': desc}


async def filter_users_parameters(
        since: Optional[str] = Query(default=None),
        desc: Optional[bool] = Query(default=False),
        limit: Optional[int] = Query(default=100, ge=1),
):
    return {'limit': limit, 'since': since, 'desc': desc}