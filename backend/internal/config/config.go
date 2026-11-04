package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config menampung seluruh konfigurasi aplikasi. Nilai dibaca dari environment
// (atau dari berkas .env sebagai pelengkap) tepat satu kali saat start.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	BaseURL       string
	AllowedOrigin string
	Port          string
}

// Load membaca .env sekali lalu memvalidasi variabel yang wajib ada.
// Variabel yang sudah disetel di environment tidak ditimpa oleh .env.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:        get("DB_HOST", "localhost"),
		DBPort:        get("DB_PORT", "5432"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     get("DB_SSLMODE", "disable"),
		BaseURL:       strings.TrimRight(os.Getenv("BASE_URL"), "/"),
		AllowedOrigin: get("ALLOWED_ORIGIN", "http://localhost:5173"),
		Port:          get("APP_PORT", "8000"),
	}

	// Fail fast: lebih baik menolak start daripada berjalan dengan nilai kosong
	// dan gagal jauh di dalam query atau response.
	var missing []string
	for _, v := range []struct{ name, value string }{
		{"DB_USER", cfg.DBUser},
		{"DB_PASSWORD", cfg.DBPassword},
		{"DB_NAME", cfg.DBName},
		{"BASE_URL", cfg.BaseURL},
	} {
		if v.value == "" {
			missing = append(missing, v.name)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf(
			"konfigurasi belum lengkap: %s — salin backend/.env.example ke backend/.env lalu isi nilainya",
			strings.Join(missing, ", "),
		)
	}

	return cfg, nil
}

func get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
