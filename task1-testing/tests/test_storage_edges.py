from app.storage import InventoryStore
from app.models import ItemCreate, ItemUpdate


def test_list_items_negative_offset_and_none_limit():
    s = InventoryStore()
    # наполняем 3 предмета
    for i in range(3):
        s.create_item(ItemCreate(
            name=f"n{i}",
            quantity=1,
            price=1.0,
            tags=[],
            status="active",
        ))

    # покрываем ветку offset < 0
    out = s.list_items(offset=-10, limit=2)
    assert len(out) == 2
    assert [x.name for x in out] == ["n0", "n1"]

    # покрываем ветку limit is None
    out2 = s.list_items(limit=None)
    assert len(out2) == 3
    assert [x.name for x in out2] == ["n0", "n1", "n2"]


def test_update_tags_and_status_branches():
    s = InventoryStore()
    it = s.create_item(ItemCreate(
        name="x",
        quantity=1,
        price=1.0,
        tags=["a"],
        status="active",
    ))

    # покрываем ветку patch.tags is not None
    it2 = s.update_item(it.id, ItemUpdate(tags=["a", "b"]))
    assert it2.tags == ["a", "b"]

    # покрываем ветку patch.status is not None
    it3 = s.update_item(it.id, ItemUpdate(status="archived"))
    assert it3.status == "archived"
