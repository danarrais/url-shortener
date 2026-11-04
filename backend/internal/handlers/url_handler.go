package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
)

const (
	// maxBodyBytes membatasi ukuran body request agar tidak bisa dipakai
	// untuk menghabiskan memori server.
	maxBodyBytes = 1 << 20 // 1 MB

	// maxURLLength adalah batas panjang URL yang wajar (batas praktis browser).
	maxURLLength = 2048
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
	// 1. Batasi ukuran body, lalu baca input JSON dari user (React)
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	var req models.ShortenRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			http.Error(w, "Body request terlalu besar", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	// 2. Validasi URL: tidak kosong, tidak kepanjangan, dan benar-benar URL
	longURL := strings.TrimSpace(req.LongURL)
	if longURL == "" || len(longURL) > maxURLLength {
		http.Error(w, "URL tidak valid", http.StatusBadRequest)
		return
	}

	parsed, err := url.ParseRequestURI(longURL)
	if err != nil || parsed.Host == "" {
		http.Error(w, "URL tidak valid", http.StatusBadRequest)
		return
	}

	// Hanya izinkan http/https. Tanpa ini, skema seperti "javascript:" bisa
	// tersimpan dan ikut terkirim sebagai header Location saat redirect.
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		http.Error(w, "Hanya URL http dan https yang diizinkan", http.StatusBadRequest)
		return
	}

	// 3. Buat kode unik 6 karakter menggunakan utils
	shortCode := utils.GenerateShortCode(6)

	// 4. Siapkan data untuk disimpan ke database (pakai URL yang sudah di-trim)
	newURL := models.URL{
		LongURL:   longURL,
		ShortCode: shortCode,
	}

	// 5. Simpan ke database melalui koki (Repository)
	err = h.Repo.Create(ctx, &newURL)
	if err != nil {
		log.Printf("ShortenURL: gagal menyimpan URL: %v", err)
		http.Error(w, "Gagal menyimpan ke database", http.StatusInternalServerError)
		return
	}

	// 6. Siapkan balasan sukses
	res := models.ShortenResponse{
		ShortURL: "http://localhost:8000/" + shortCode,
	}

	// 7. Kirim balasan JSON ke user
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
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "URL tidak ditemukan", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Redirect: gagal mencari short code %q: %v", shortCode, err)
		http.Error(w, "Terjadi kesalahan pada server", http.StatusInternalServerError)
		return
	}

	// 3. Lempar (Redirect) user ke URL asli
	_ = h.Repo.IncrementClicks(ctx, shortCode)

	// 302 (Found), bukan 301: 301 di-cache permanen oleh browser, sehingga
	// klik berikutnya tidak sampai ke server dan statistik tidak akurat.
	http.Redirect(w, r, urlData.LongURL, http.StatusFound)
}

func (h *URLHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortCode := r.PathValue("short_code")

	urlData, err := h.Repo.FindByShortCode(ctx, shortCode)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "URL tidak ditemukan", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("GetStats: gagal mencari short code %q: %v", shortCode, err)
		http.Error(w, "Terjadi kesalahan pada server", http.StatusInternalServerError)
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
