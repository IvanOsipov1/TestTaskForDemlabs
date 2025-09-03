import math

from app.models import Item
from app.storage import calculate_total_value


def test_calculate_total_value_basic():
    items = [
        Item(id=1, name="a", quantity=2, price=3.5, tags=[], status="active"),
        Item(id=2, name="b", quantity=0, price=10.0, tags=[], status="archived"),
        Item(id=3, name="c", quantity=5, price=1.2, tags=["x"], status="active"),
    ]
    total = calculate_total_value(items)
    # 2*3.5 + 0*10 + 5*1.2 = 7.0 + 0 + 6.0 = 13.0
    assert math.isclose(total, 13.0, rel_tol=1e-9)

def test_calculate_total_value_large():
    items = [
        Item(id=1, name="x", quantity=1000000, price=0.01, tags=[], status="active"),
        Item(id=2, name="y", quantity=1, price=9999.99, tags=[], status="active"),
    ]
    total = calculate_total_value(items)
    # 1_000_000 * 0.01 + 1 * 9999.99 = 10000 + 9999.99 = 19999.99
    assert math.isclose(total, 19999.99, rel_tol=1e-9)
