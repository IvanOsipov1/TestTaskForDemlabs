package apihttp

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/yourname/inventory-go/internal/model"
	"github.com/yourname/inventory-go/internal/store"
)

func newTestServer(t *testing.T) (*httptest.Server, *store.Store) {
	t.Helper()
	db, err := store.Open("file::memory:?cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatalf("open mem db: %v", err)
	}
	st := store.New(db)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := st.Migrate(ctx); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	ts := httptest.NewServer(Router(st))
	t.Cleanup(func() {
		ts.Close()
		db.Close()
	})
	return ts, st
}

func postJSON(t *testing.T, client *http.Client, url string, body any) (*http.Response, map[string]any) {
	t.Helper()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	var m map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&m)
	resp.Body.Close()
	return resp, m
}

func patchJSON(t *testing.T, client *http.Client, url string, body any) (*http.Response, map[string]any) {
	t.Helper()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("PATCH %s: %v", url, err)
	}
	var m map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&m)
	resp.Body.Close()
	return resp, m
}

func getJSON(t *testing.T, client *http.Client, url string) (*http.Response, map[string]any) {
	t.Helper()
	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	var m map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&m)
	resp.Body.Close()
	return resp, m
}

func deleteReq(t *testing.T, client *http.Client, url string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("DELETE %s: %v", url, err)
	}
	resp.Body.Close()
	return resp
}

func TestPOST_GET_HappyPath(t *testing.T) {
	ts, _ := newTestServer(t)
	client := ts.Client()

	resp, m := postJSON(t, client, ts.URL+"/items", model.ItemCreate{
		Name:     "pen",
		Quantity: 10,
		Price:    1.5,
		Tags:     []string{"stationery"},
		Status:   "active",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %v", resp.StatusCode, m)
	}
	id := int(m["id"].(float64))

	resp2, m2 := getJSON(t, client, ts.URL+"/items/"+chi.URLParamFromCtx(withID(t, id), "id"))
	_ = resp2
	_ = m2

	// Прямой GET по URL (без chi helper)
	resp3, m3 := getJSON(t, client, ts.URL+"/items/"+itoa(id))
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d: %v", resp3.StatusCode, m3)
	}
	if m3["name"] != "pen" {
		t.Fatalf("expected name=pen, got %v", m3["name"])
	}
}

func TestDuplicateNameConflict(t *testing.T) {
	ts, _ := newTestServer(t)
	client := ts.Client()

	resp, _ := postJSON(t, client, ts.URL+"/items", model.ItemCreate{
		Name:     "dupe",
		Quantity: 1,
		Price:    2.0,
		Tags:     []string{},
		Status:   "active",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	resp2, _ := postJSON(t, client, ts.URL+"/items", model.ItemCreate{
		Name:     "dupe",
		Quantity: 2,
		Price:    3.0,
		Tags:     []string{},
		Status:   "active",
	})
	if resp2.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", resp2.StatusCode)
	}
}

func TestGetNotFound(t *testing.T) {
	ts, _ := newTestServer(t)
	client := ts.Client()

	resp, _ := getJSON(t, client, ts.URL+"/items/99999")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestPatchRenameConflict(t *testing.T) {
	ts, _ := newTestServer(t)
	client := ts.Client()

	_, a := postJSON(t, client, ts.URL+"/items", model.ItemCreate{
		Name:     "a",
		Quantity: 1,
		Price:    1.0,
		Tags:     []string{},
		Status:   "active",
	})
	_, b := postJSON(t, client, ts.URL+"/items", model.ItemCreate{
		Name:     "b",
		Quantity: 1,
		Price:    1.0,
		Tags:     []string{},
		Status:   "active",
	})

	idA := int(a["id"].(float64))
	_ = b

	resp, body := patchJSON(t, client, ts.URL+"/items/"+itoa(idA), map[string]any{"name": "b"})
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d body=%v", resp.StatusCode, body)
	}
}

func TestValidation422(t *testing.T) {
	ts, _ := newTestServer(t)
	client := ts.Client()

	resp, _ := postJSON(t, client, ts.URL+"/items", map[string]any{
		"name":     "",
		"quantity": -1,
		"price":    0,
		"tags":     []string{},
		"status":   "unknown",
	})
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", resp.StatusCode)
	}
}

func TestListLimitOffsetAndDelete(t *testing.T) {
	ts, _ := newTestServer(t)
	client := ts.Client()

	ids := make([]int, 0, 3)
	for i := 0; i < 3; i++ {
		resp, m := postJSON(t, client, ts.URL+"/items", model.ItemCreate{
			Name:     "n" + itoa(i),
			Quantity: i,
			Price:    1.0,
			Tags:     []string{},
			Status:   "active",
		})
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201, got %d", resp.StatusCode)
		}
		ids = append(ids, int(m["id"].(float64)))
	}

	resp, m := getJSON(t, client, ts.URL+"/items?limit=2&offset=1")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	arr, ok := m[""].([]any) // intentionally wrong to check structure
	if ok && len(arr) > 0 {
		t.Fatalf("unexpected structure")
	}

	// Прочитаем список как []map
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/items?limit=2&offset=1", nil)
	resp2, err := client.Do(req)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var list []map[string]any
	_ = json.NewDecoder(resp2.Body).Decode(&list)
	resp2.Body.Close()
	if len(list) != 2 {
		t.Fatalf("expected 2, got %d", len(list))
	}

	// DELETE и проверка 404
	delResp := deleteReq(t, client, ts.URL+"/items/"+itoa(ids[0]))
	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", delResp.StatusCode)
	}
	resp404, _ := getJSON(t, client, ts.URL+"/items/"+itoa(ids[0]))
	if resp404.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", resp404.StatusCode)
	}
}

// helpers

func itoa(i int) string {
	return strconvItoa(i)
}

// avoid shadowing strconv in imports above
func strconvItoa(i int) string {
	return strconvItoaImpl(i)
}

func strconvItoaImpl(i int) string {
	return strconvItoaStd(i)
}

// isolate std import here to not clutter above
func strconvItoaStd(i int) string {
	return strconv.Itoa(i)
}

func withID(t *testing.T, id int) *http.Request {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, "/items/"+itoa(id), nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", itoa(id))
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
