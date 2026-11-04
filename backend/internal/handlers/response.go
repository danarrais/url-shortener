package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// Semua response API berbentuk JSON. Bentuk error selalu sama:
//
//	{"error": "pesan"}
//
// Frontend (readError di frontend/src/App.jsx) membaca field "error" tersebut.

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Header sudah terkirim, jadi status tidak bisa diubah lagi.
		// Cukup catat supaya tidak hilang diam-diam.
		log.Printf("writeJSON: gagal menulis response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
