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

	// 1. Koneksi Database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Berhasil terhubung ke Database PostgreSQL!")

	// 2. Auto-Migrate Tabel
	err = db.AutoMigrate(&models.URL{})
	if err != nil {
		log.Fatal("Gagal menemukan migrasi database", err)
	}
	fmt.Println("Migrasi tabel 'urls' berhasil!")

	// 3. DEPENDENCY INJECTION (Menyambungkan komponen)
	// Kita suntikkan 'db' ke dalam koki (Repo)
	urlRepo := repository.NewURLRepository(db)

	// Lalu koki (Repo) disuntikkan ke dalam pelayan (Handler)
	urlHandler := handlers.NewURLHandler(urlRepo)

	// 4. Siapkan Router (Mux)
	mux := http.NewServeMux()

	// 5. Daftarkan Rute API
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Server URL Shortener siap menerima perintah!")
	})

	// Rute POST untuk memendekkan URL
	mux.HandleFunc("POST /shorten", urlHandler.ShortenURL)

	mux.HandleFunc("GET /stats/{short_code}", urlHandler.GetStats)
	
	// Rute GET untuk redirect (menggunakan fitur Path Parameter Go 1.22)
	mux.HandleFunc("GET /{short_code}", urlHandler.Redirect)

	// 6. Nyalakan Server
	fmt.Println("Server menyala di http://localhost:8000")
	err = http.ListenAndServe(":8000", enableCORS(mux))
	if err != nil {
		log.Fatal("Gagal menyalakan server:", err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Memberikan izin masuk kepada aplikasi React
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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
