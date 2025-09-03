import pytest
from fastapi.testclient import TestClient

from app.main import app, store
from app.models import ItemCreate


@pytest.fixture(autouse=True)
def _reset_store():
    store.reset()
    yield
    store.reset()


@pytest.fixture
def client():
    return TestClient(app)


def _preload_2k():
    # Заполняем напрямую стор (быстрее), бенчмарк — только GET
    for i in range(2000):
        payload = ItemCreate(
            name=f"item-{i}",
            quantity=(i % 10) + 1,
            price=1.0 + (i % 7) * 0.1,
            tags=["bench", str(i % 5)],
            status="active",
        )
        store.create_item(payload)


def test_benchmark_get_items_first_500(client: TestClient, benchmark):
    _preload_2k()

    def fetch():
        r = client.get("/items", params={"limit": 500, "offset": 0})
        assert r.status_code == 200
        data = r.json()
        assert len(data) == 500
        return data

    # Запускаем бенчмарк только для GET
    benchmark(fetch)
