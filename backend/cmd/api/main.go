package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/repository"
)

func main() {
	// 1. Baca seluruh konfigurasi sekali di awal
	cfg, err := config.Load()
	if err != nil {
		slog.Error("konfigurasi tidak valid", "error", err)
		os.Exit(1)
	}

	// 2. Koneksi Database
	db, err := config.ConnectDB(cfg)
	if err != nil {
		slog.Error("gagal terhubung ke database", "error", err)
		os.Exit(1)
	}
	slog.Info("terhubung ke database", "host", cfg.DBHost, "database", cfg.DBName)

	// 3. Auto-Migrate Tabel
	if err := db.AutoMigrate(&models.URL{}); err != nil {
		slog.Error("gagal menjalankan migrasi database", "error", err)
		os.Exit(1)
	}
	slog.Info("migrasi tabel berhasil", "tabel", "urls")

	// 4. DEPENDENCY INJECTION (Menyambungkan komponen)
	// Kita suntikkan 'db' ke dalam koki (Repo)
	urlRepo := repository.NewURLRepository(db)

	// Lalu koki (Repo) disuntikkan ke dalam pelayan (Handler)
	urlHandler := handlers.NewURLHandler(urlRepo, cfg.BaseURL)

	// 5. Siapkan Router (Mux)
	mux := http.NewServeMux()

	// 6. Daftarkan Rute API
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Server URL Shortener siap menerima perintah!")
	})

	// Rute POST untuk memendekkan URL
	mux.HandleFunc("POST /shorten", urlHandler.ShortenURL)

	mux.HandleFunc("GET /stats/{short_code}", urlHandler.GetStats)

	// Rute GET untuk redirect (menggunakan fitur Path Parameter Go 1.22)
	mux.HandleFunc("GET /{short_code}", urlHandler.Redirect)

	// 7. Server dengan batas waktu: tanpa timeout, koneksi yang menggantung
	// (slowloris) bisa menahan goroutine dan memori server tanpa batas.
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           middleware.CORS(cfg.AllowedOrigin)(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Ctrl+C atau SIGTERM membatalkan ctx, bukan mematikan proses seketika.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server menyala", "alamat", cfg.BaseURL)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server berhenti karena error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("sinyal berhenti diterima, menutup server")

	// Beri request yang sedang berjalan waktu 10 detik untuk selesai.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("gagal menutup server dengan rapi", "error", err)
		os.Exit(1)
	}
	slog.Info("server berhenti dengan rapi")
}
