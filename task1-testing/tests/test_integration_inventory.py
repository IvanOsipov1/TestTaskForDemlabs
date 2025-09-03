import pytest
from fastapi.testclient import TestClient

from app.main import app, store


@pytest.fixture(autouse=True)
def _reset_store():
    store.reset()
    yield
    store.reset()


@pytest.fixture
def client():
    return TestClient(app)


def _mk_item_payload(name="pen", quantity=10, price=1.5, tags=None, status="active"):
    return {
        "name": name,
        "quantity": quantity,
        "price": price,
        "tags": tags or ["default"],
        "status": status,
    }


def test_create_and_get_item(client: TestClient):
    r = client.post("/items", json=_mk_item_payload())
    assert r.status_code == 201, r.text
    item = r.json()
    assert item["id"] > 0
    assert item["name"] == "pen"

    rid = item["id"]
    r2 = client.get(f"/items/{rid}")
    assert r2.status_code == 200
    assert r2.json()["name"] == "pen"


def test_list_items_with_limit(client: TestClient):
    for i in range(5):
        payload = _mk_item_payload(name=f"itm{i}", quantity=i, price=1.0 + i)
        assert client.post("/items", json=payload).status_code == 201
    r = client.get("/items", params={"limit": 2, "offset": 1})
    assert r.status_code == 200
    data = r.json()
    assert len(data) == 2
    assert data[0]["name"] == "itm1"
    assert data[1]["name"] == "itm2"


def test_patch_update_and_rename(client: TestClient):
    a = client.post("/items", json=_mk_item_payload(name="a")).json()
    b = client.post("/items", json=_mk_item_payload(name="b")).json()

    # update quantity & price
    r = client.patch(f"/items/{a['id']}", json={"quantity": 123, "price": 9.99})
    assert r.status_code == 200
    updated = r.json()
    assert updated["quantity"] == 123
    assert updated["price"] == 9.99

    # rename to non-conflicting name
    r2 = client.patch(f"/items/{a['id']}", json={"name": "a2"})
    assert r2.status_code == 200
    assert r2.json()["name"] == "a2"

    # try rename to existing name -> 409
    r3 = client.patch(f"/items/{a['id']}", json={"name": "b"})
    assert r3.status_code == 409
    assert r3.json()["detail"] == "name already exists"


def test_delete_then_404(client: TestClient):
    item = client.post("/items", json=_mk_item_payload(name="delme")).json()
    rid = item["id"]

    r = client.delete(f"/items/{rid}")
    assert r.status_code == 204

    r2 = client.get(f"/items/{rid}")
    assert r2.status_code == 404
    assert r2.json()["detail"] == "item not found"


def test_duplicate_name_on_create(client: TestClient):
    assert client.post("/items", json=_mk_item_payload(name="uniq")).status_code == 201
    r = client.post("/items", json=_mk_item_payload(name="uniq"))
    assert r.status_code == 409
    assert r.json()["detail"] == "name already exists"


def test_404s(client: TestClient):
    assert client.get("/items/9999").status_code == 404
    assert client.patch("/items/9999", json={"name": "x"}).status_code == 404
    assert client.delete("/items/9999").status_code == 404


def test_invalid_data_422(client: TestClient):
    # empty name
    assert client.post("/items", json=_mk_item_payload(name="")).status_code == 422
    # negative quantity
    assert client.post("/items", json=_mk_item_payload(quantity=-1)).status_code == 422
    # zero/non-positive price
    assert client.post("/items", json=_mk_item_payload(price=0)).status_code == 422
    # invalid status
    bad = _mk_item_payload(status="unknown")
    r = client.post("/items", json=bad)
    assert r.status_code == 422


def test_stats_total_value(client: TestClient):
    client.post("/items", json=_mk_item_payload(name="x", quantity=2, price=3.5))
    client.post("/items", json=_mk_item_payload(name="y", quantity=5, price=1.2))
    r = client.get("/stats/total_value")
    assert r.status_code == 200
    data = r.json()
    assert abs(data["total_value"] - 13.0) < 1e-9
