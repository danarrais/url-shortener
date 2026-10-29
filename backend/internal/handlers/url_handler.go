package handlers

import(
	"encoding/json"
	"net/http"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
)

// URLHandler menyimpan 'koki' (repository) agar bisa menyuruhnya memasak data
type URLHandler struct {
	Repo *repository.URLRepository
}

// NewURLHandler adalah fungsi untuk membuat pelayan baru
func NewURLHandler(repo *repository.URLRepository) *URLHandler {
	return &URLHandler{Repo: repo}
}

// ShortenURL bertugas menerima URL panjang dan membalas URL pendek (POST)
func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	// 1. Baca input JSON dari user (React)
	ctx := r.Context()

	var req models.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	// 2. Buat kode unik 6 karakter menggunakan utils
	shortCode := utils.GenerateShortCode(6)

	// 3. Siapkan data untuk disimpan ke database
	newURL := models.URL{
		LongURL: req.LongURL,
		ShortCode: shortCode,
	}

	// 4. Simpan ke database melalui koki (Repository)
	err = h.Repo.Create(ctx, &newURL)
	if err != nil {
		http.Error(w, "Gagal menyimpan ke database", http.StatusInternalServerError)
		return
	}

	// 5. Siapkan balasan sukses
	res := models.ShortenResponse {
		ShortURL: "http://localhost:8000/" + shortCode,
	}

	// 6. Kirim balasan JSON ke user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// Redirect bertugas menangkap klik URL pendek dan melempar user ke URL asli (GET)
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil kode dari ujung URL
	ctx := r.Context()

	shortCode := r.PathValue("short_code")

	// 2. Cari data aslinya di database
	urlData, err := h.Repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		http.Error(w, "URL tidak ditemukan", http.StatusNotFound)
		return
	}

	// 3. Lempar (Redirect) user ke URL asli
	_ = h.Repo.IncrementClicks(ctx, shortCode)

	http.Redirect(w, r, urlData.LongURL, http.StatusFound)
}

func (h *URLHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortCode := r.PathValue("short_code")

	urlData, err := h.Repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		http.Error(w, "URL tidak ditemukan", http.StatusNotFound)
		return
	}

	stats := map[string]interface{}{
		"short_code": urlData.ShortCode,
		"long_url":   urlData.LongURL,
		"clicks":     urlData.Clicks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}