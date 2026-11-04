package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/models"
	"backend/internal/repository"
)

func main() {
	fmt.Println("Memulai aplikasi...")

	// 1. Baca seluruh konfigurasi sekali di awal
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Koneksi Database
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Berhasil terhubung ke Database PostgreSQL!")

	// 3. Auto-Migrate Tabel
	err = db.AutoMigrate(&models.URL{})
	if err != nil {
		log.Fatal("Gagal menjalankan migrasi database: ", err)
	}
	fmt.Println("Migrasi tabel 'urls' berhasil!")

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

	// 7. Nyalakan Server
	fmt.Printf("Server menyala di %s\n", cfg.BaseURL)
	err = http.ListenAndServe(":"+cfg.Port, enableCORS(mux, cfg.AllowedOrigin))
	if err != nil {
		log.Fatal("Gagal menyalakan server:", err)
	}
}

func enableCORS(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Memberikan izin masuk kepada aplikasi React
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Menangani "preflight request" (ketukan pintu awal dari browser)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
