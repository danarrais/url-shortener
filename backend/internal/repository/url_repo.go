package repository

import (
	"backend/internal/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ErrNotFound dikembalikan ketika short code tidak ada di database.
// Handler memakai error ini untuk membedakan 404 dari kegagalan server (500).
var ErrNotFound = errors.New("url tidak ditemukan")

type URLRepository struct {
	DB *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{DB: db}
}

func (r *URLRepository) Create(ctx context.Context, url *models.URL) error {
	return r.DB.Create(url).Error
}

func (r *URLRepository) FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	var url models.URL

	err := r.DB.Where("short_code = ? ", shortCode).First(&url).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find by short code: %w", err)
	}

	return &url, nil
}

func (r *URLRepository) IncrementClicks(ctx context.Context, shortCode string) error {
	return r.DB.WithContext(ctx).Model(&models.URL{}).
		Where("short_code = ?", shortCode).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", 1)).Error
}