from http import HTTPStatus
from typing import Optional, Dict, Union, Any

from fastapi import HTTPException

# TODO передавать detail в конструкторе


class HttpNotFoundException(HTTPException):
    status_code = HTTPStatus.NOT_FOUND
    detail = 'Item not found'

    def __init__(self):
        super().__init__(status_code=self.status_code, detail=self.detail)

    @classmethod
    def responses_dict(cls) -> Optional[Dict[Union[int, str], Dict[str, Any]]]:
        return {
            'description': cls.detail,
            'content': {'application/json': {'example': {'message': cls.detail}}}
        }


class HttpConflictException(HTTPException):
    status_code = HTTPStatus.CONFLICT
    detail = 'Item must be unique'

    def __init__(self):
        super().__init__(status_code=self.status_code, detail=self.detail)

    @classmethod
    def responses_dict(cls) -> Optional[Dict[Union[int, str], Dict[str, Any]]]:
        return {
            'description': cls.detail,
            'content': {'application/json': {'example': {'message': cls.detail}}}
        }
