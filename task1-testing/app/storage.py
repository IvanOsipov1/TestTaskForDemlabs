from __future__ import annotations

from threading import Lock
from typing import Dict, List, Optional, Tuple

from .models import Item, ItemCreate, ItemUpdate


class NotFoundError(Exception):
    pass


class DuplicateNameError(Exception):
    pass


class InventoryStore:
    """
    Простое in-memory хранилище с индексом уникальных имён.
    Потокобезопасность обеспечена через Lock (на случай продвинутых сценариев).
    """

    def __init__(self) -> None:
        self._lock = Lock()
        self._items: Dict[int, Item] = {}
        self.name_index: Dict[str, int] = {}
        self._id_seq: int = 1

    # --- service / testing ---
    def reset(self) -> None:
        with self._lock:
            self._items.clear()
            self.name_index.clear()
            self._id_seq = 1

    # --- core operations ---
    def create_item(self, data: ItemCreate) -> Item:
        with self._lock:
            if data.name in self.name_index:
                raise DuplicateNameError(data.name)

            item_id = self._id_seq
            self._id_seq += 1

            item = Item(id=item_id, **data.model_dump())
            self._items[item_id] = item
            self.name_index[data.name] = item_id
            return item

    def get_item(self, item_id: int) -> Item:
        with self._lock:
            item = self._items.get(item_id)
            if item is None:
                raise NotFoundError(item_id)
            return item

    def list_items(self, offset: int = 0, limit: Optional[int] = None) -> List[Item]:
        with self._lock:
            items = list(sorted(self._items.values(), key=lambda x: x.id))
            if offset < 0:
                offset = 0
            if limit is None:
                return items[offset:]
            return items[offset: offset + limit]

    def delete_item(self, item_id: int) -> None:
        with self._lock:
            item = self._items.pop(item_id, None)
            if item is None:
                raise NotFoundError(item_id)
            # удалить из индекса
            self.name_index.pop(item.name, None)

    def update_item(self, item_id: int, patch: ItemUpdate) -> Item:
        with self._lock:
            item = self._items.get(item_id)
            if item is None:
                raise NotFoundError(item_id)

            # обработка смены имени с проверкой конфликтов
            if patch.name is not None and patch.name != item.name:
                if patch.name in self.name_index and self.name_index[patch.name] != item_id:
                    raise DuplicateNameError(patch.name)
                # переиндексация
                old_name = item.name
                item.name = patch.name
                self.name_index.pop(old_name, None)
                self.name_index[item.name] = item_id

            if patch.quantity is not None:
                item.quantity = patch.quantity

            if patch.price is not None:
                item.price = patch.price

            if patch.tags is not None:
                item.tags = patch.tags

            if patch.status is not None:
                item.status = patch.status

            self._items[item_id] = item
            return item


def calculate_total_value(items: List[Item]) -> float:
    """
    Чистая бизнес-логика: сумма price * quantity по списку Item.
    """
    total = 0.0
    for it in items:
        total += float(it.price) * int(it.quantity)
    return total
