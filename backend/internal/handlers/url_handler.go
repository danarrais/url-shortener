package handlers

import (
	"context"
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

// URLStore adalah kontrak data yang dibutuhkan handler. Interface ditulis di
// package handlers agar handler tidak terikat pada GORM dan bisa diuji dengan
// implementasi palsu tanpa database.
type URLStore interface {
	Create(ctx context.Context, url *models.URL) error
	FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error)
	IncrementClicks(ctx context.Context, shortCode string) error
}

// URLHandler menyimpan 'koki' (repository) agar bisa menyuruhnya memasak data
type URLHandler struct {
	Repo    URLStore
	BaseURL string
}

// NewURLHandler adalah fungsi untuk membuat pelayan baru
func NewURLHandler(repo URLStore, baseURL string) *URLHandler {
	return &URLHandler{Repo: repo, BaseURL: baseURL}
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
			writeError(w, http.StatusRequestEntityTooLarge, "Body request terlalu besar")
			return
		}
		writeError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	// 2. Validasi URL: tidak kosong, tidak kepanjangan, dan benar-benar URL
	longURL := strings.TrimSpace(req.LongURL)
	if longURL == "" || len(longURL) > maxURLLength {
		writeError(w, http.StatusBadRequest, "URL tidak valid")
		return
	}

	parsed, err := url.ParseRequestURI(longURL)
	if err != nil || parsed.Host == "" {
		writeError(w, http.StatusBadRequest, "URL tidak valid")
		return
	}

	// Hanya izinkan http/https. Tanpa ini, skema seperti "javascript:" bisa
	// tersimpan dan ikut terkirim sebagai header Location saat redirect.
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		writeError(w, http.StatusBadRequest, "Hanya URL http dan https yang diizinkan")
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
		writeError(w, http.StatusInternalServerError, "Gagal menyimpan ke database")
		return
	}

	// 6. Siapkan balasan sukses
	res := models.ShortenResponse{
		ShortURL: h.BaseURL + "/" + shortCode,
	}

	// 7. Kirim balasan JSON ke user
	writeJSON(w, http.StatusOK, res)
}

// Redirect bertugas menangkap klik URL pendek dan melempar user ke URL asli (GET)
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil kode dari ujung URL
	ctx := r.Context()

	shortCode := r.PathValue("short_code")

	// 2. Cari data aslinya di database
	urlData, err := h.Repo.FindByShortCode(ctx, shortCode)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(w, http.StatusNotFound, "URL tidak ditemukan")
		return
	}
	if err != nil {
		log.Printf("Redirect: gagal mencari short code %q: %v", shortCode, err)
		writeError(w, http.StatusInternalServerError, "Terjadi kesalahan pada server")
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
		writeError(w, http.StatusNotFound, "URL tidak ditemukan")
		return
	}
	if err != nil {
		log.Printf("GetStats: gagal mencari short code %q: %v", shortCode, err)
		writeError(w, http.StatusInternalServerError, "Terjadi kesalahan pada server")
		return
	}

	stats := models.StatsResponse{
		ShortCode: urlData.ShortCode,
		LongURL:   urlData.LongURL,
		Clicks:    urlData.Clicks,
	}

	writeJSON(w, http.StatusOK, stats)
}
