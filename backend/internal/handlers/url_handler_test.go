package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/models"
	"backend/internal/repository"
)

// fakeStore menggantikan repository.URLRepository agar handler dapat diuji
// tanpa database. Field disetel per test sesuai skenario.
type fakeStore struct {
	created   *models.URL
	createErr error
	data      *models.URL
	findErr   error
	clicked   string
}

func (f *fakeStore) Create(_ context.Context, url *models.URL) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = url
	return nil
}

func (f *fakeStore) FindByShortCode(_ context.Context, _ string) (*models.URL, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.data, nil
}

func (f *fakeStore) IncrementClicks(_ context.Context, shortCode string) error {
	f.clicked = shortCode
	return nil
}

// assertError memastikan response memakai bentuk error JSON yang seragam.
func assertError(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()

	if rec.Code != want {
		t.Fatalf("status = %d, mau %d", rec.Code, want)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("Content-Type = %q, mau application/json", ct)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body bukan JSON: %v", err)
	}
	if body["error"] == "" {
		t.Fatalf("field error kosong: %#v", body)
	}
}

func TestShortenURLSuccess(t *testing.T) {
	store := &fakeStore{}
	handler := NewURLHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"long_url":"  https://example.com/halaman  "}`))
	rec := httptest.NewRecorder()

	handler.ShortenURL(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, mau 200; body = %s", rec.Code, rec.Body.String())
	}

	var res models.ShortenResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("response bukan JSON: %v", err)
	}
	if res.ShortURL == "" {
		t.Fatal("short_url kosong")
	}
	if store.created == nil {
		t.Fatal("repository tidak pernah menerima data")
	}
	if want := "https://example.com/halaman"; store.created.LongURL != want {
		t.Fatalf("long_url tersimpan = %q, mau %q (sudah di-trim)", store.created.LongURL, want)
	}
}

func TestShortenURLInvalidJSON(t *testing.T) {
	handler := NewURLHandler(&fakeStore{})

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"long_url":`))
	rec := httptest.NewRecorder()

	handler.ShortenURL(rec, req)

	assertError(t, rec, http.StatusBadRequest)
}

func TestShortenURLCreateError(t *testing.T) {
	store := &fakeStore{createErr: errors.New("koneksi database putus")}
	handler := NewURLHandler(store)

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"long_url":"https://example.com"}`))
	rec := httptest.NewRecorder()

	handler.ShortenURL(rec, req)

	assertError(t, rec, http.StatusInternalServerError)
}

func TestRedirectFound(t *testing.T) {
	store := &fakeStore{data: &models.URL{ShortCode: "abc123", LongURL: "https://example.com/tujuan"}}
	handler := NewURLHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	req.SetPathValue("short_code", "abc123")
	rec := httptest.NewRecorder()

	handler.Redirect(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, mau 302", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "https://example.com/tujuan" {
		t.Fatalf("Location = %q, mau %q", loc, "https://example.com/tujuan")
	}
	if store.clicked != "abc123" {
		t.Fatalf("IncrementClicks dipanggil dengan %q, mau %q", store.clicked, "abc123")
	}
}

func TestRedirectNotFound(t *testing.T) {
	handler := NewURLHandler(&fakeStore{findErr: repository.ErrNotFound})

	req := httptest.NewRequest(http.MethodGet, "/zzzzzz", nil)
	req.SetPathValue("short_code", "zzzzzz")
	rec := httptest.NewRecorder()

	handler.Redirect(rec, req)

	assertError(t, rec, http.StatusNotFound)
}

func TestGetStatsOK(t *testing.T) {
	store := &fakeStore{data: &models.URL{ShortCode: "abc123", LongURL: "https://example.com", Clicks: 7}}
	handler := NewURLHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/stats/abc123", nil)
	req.SetPathValue("short_code", "abc123")
	rec := httptest.NewRecorder()

	handler.GetStats(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, mau 200; body = %s", rec.Code, rec.Body.String())
	}

	var got models.StatsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response bukan JSON: %v", err)
	}
	if got.Clicks != 7 || got.ShortCode != "abc123" {
		t.Fatalf("stats = %+v, mau short_code=abc123 clicks=7", got)
	}
}
