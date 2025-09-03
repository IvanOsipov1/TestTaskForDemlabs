from typing import List, Literal, Optional
from pydantic import BaseModel, Field

StatusLiteral = Literal["active", "archived"]


class ItemBase(BaseModel):
    name: str = Field(..., min_length=1, max_length=64)
    quantity: int = Field(..., ge=0)
    price: float = Field(..., gt=0)
    tags: List[str] = Field(default_factory=list)
    status: StatusLiteral


class ItemCreate(ItemBase):
    """Payload для POST /items."""


class ItemUpdate(BaseModel):
    """Частичное обновление для PATCH /items/{id}."""
    name: Optional[str] = Field(None, min_length=1, max_length=64)
    quantity: Optional[int] = Field(None, ge=0)
    price: Optional[float] = Field(None, gt=0)
    tags: Optional[List[str]] = None
    status: Optional[StatusLiteral] = None


class Item(ItemBase):
    id: int
