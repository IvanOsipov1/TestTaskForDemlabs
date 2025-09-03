from typing import List, Optional

from fastapi import FastAPI, HTTPException, Query, status

from .models import Item, ItemCreate, ItemUpdate
from .storage import InventoryStore, NotFoundError, DuplicateNameError, calculate_total_value

app = FastAPI(title="Inventory (in-memory)")

# Глобальный стор — общий для приложения и тестов
store = InventoryStore()


@app.post("/items", response_model=Item, status_code=status.HTTP_201_CREATED)
def create_item(payload: ItemCreate) -> Item:
    try:
        return store.create_item(payload)
    except DuplicateNameError:
        raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="name already exists")


@app.get("/items", response_model=List[Item])
def list_items(
    limit: Optional[int] = Query(None, ge=1, le=10000),
    offset: int = Query(0, ge=0),
) -> List[Item]:
    return store.list_items(offset=offset, limit=limit)


@app.get("/items/{item_id}", response_model=Item)
def get_item(item_id: int) -> Item:
    try:
        return store.get_item(item_id)
    except NotFoundError:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="item not found")


@app.patch("/items/{item_id}", response_model=Item)
def patch_item(item_id: int, payload: ItemUpdate) -> Item:
    try:
        return store.update_item(item_id, payload)
    except NotFoundError:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="item not found")
    except DuplicateNameError:
        raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="name already exists")


@app.delete("/items/{item_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_item(item_id: int) -> None:
    try:
        store.delete_item(item_id)
    except NotFoundError:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="item not found")


@app.get("/stats/total_value")
def stats_total_value() -> dict:
    items = store.list_items()
    return {"total_value": calculate_total_value(items)}
