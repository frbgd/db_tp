import logging

from fastapi import FastAPI, Request
import uvicorn as uvicorn
from fastapi.responses import ORJSONResponse

from api import forum, post, service, thread, user
from core.config import dsl
from core.logger import LOGGING
from db import postgres
from db.postgres import PostgressDbEngine

app = FastAPI(
    default_response_class=ORJSONResponse,
)


app.include_router(forum.router, prefix='/api/forum')
app.include_router(post.router, prefix='/api/post')
app.include_router(service.router, prefix='/api/service')
app.include_router(thread.router, prefix='/api/thread')
app.include_router(user.router, prefix='/api/user')


@app.on_event("startup")
async def startup_event():
    postgres.postgres_db_engine = PostgressDbEngine(**dsl)
    await postgres.postgres_db_engine.connect()


@app.on_event("shutdown")
async def shutdown_event():
    await postgres.postgres_db_engine.close()


if __name__ == '__main__':
    uvicorn.run(
        'main:app',
        host='0.0.0.0',
        port=5000,
        log_config=LOGGING,
        log_level=logging.DEBUG,
    )
