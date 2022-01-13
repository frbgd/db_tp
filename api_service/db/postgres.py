import asyncpg

from typing import Dict, Any

from asyncpg import Connection, UniqueViolationError, NotNullViolationError, ForeignKeyViolationError


async def create_postgres_instance(**dsl: Dict[str, Any]) -> Connection:
    """Потому что в конструкторе класса нельзя вызывать await"""
    return await asyncpg.connect(**dsl)


item_not_unique_exception = UniqueViolationError
not_null_constraint_exception = NotNullViolationError
foreign_key_violation_exception = ForeignKeyViolationError


class PostgressDbEngine:
    def __init__(self, user, password, host, database, port="5432"):
        self.user = user
        self.password = password
        self.host = host
        self.port = port
        self.database = database
        self._cursor = None

        self._connection_pool = None
        self.con = None

    async def connect(self):
        if not self._connection_pool:
            self._connection_pool = await asyncpg.create_pool(
                min_size=1,
                max_size=5000,
                command_timeout=60,
                host=self.host,
                port=self.port,
                user=self.user,
                password=self.password,
                database=self.database
            )

    async def close(self):
        if not self._connection_pool:
            await self._connection_pool.close()

    async def _fetchrow(self, sql: str, *args):
        if not self._connection_pool:
            await self.connect()
        con = await self._connection_pool.acquire()
        try:
            result = await con.fetchrow(sql, *args)
            return result
        except Exception as e:
            raise
        finally:
            await self._connection_pool.release(con)

    async def _fetch(self, sql: str, *args):
        if not self._connection_pool:
            await self.connect()
        con = await self._connection_pool.acquire()
        try:
            result = await con.fetch(sql, *args)
            return result
        except Exception as e:
            raise
        finally:
            await self._connection_pool.release(con)

    async def get_one(self, sql: str, *args):
        return await self._fetchrow(sql, *args)

    async def get_many(self, sql: str, *args):
        return await self._fetch(sql, *args)

    async def insert(self, sql: str, *args):
        return await self._fetchrow(sql, *args)

    async def insert_many(self, sql: str, *args):
        return await self._fetch(sql, *args)

    async def update(self, sql: str, *args):
        return await self._fetchrow(sql, *args)

    async def execute(self, sql, *args):
        if not self._connection_pool:
            await self.connect()
        con = await self._connection_pool.acquire()
        try:
            result = await con.execute(sql, *args)
            return result
        except Exception as e:
            raise
        finally:
            await self._connection_pool.release(con)


postgres_db_engine: PostgressDbEngine = None


async def get_postgres() -> PostgressDbEngine:
    return postgres_db_engine
