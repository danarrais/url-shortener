package repository

import (
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type URLRepository struct {
	DB *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{DB: db}
}

func (r *URLRepository) Create(ctx context.Context, url *models.URL) error {
	return r.DB.Create(url).Error
}

func (r *URLRepository) FindByShortCode(ctx context.Context,shortCode string) (*models.URL, error) {
	var url models.URL

	err := r.DB.Where("short_code = ? ", shortCode).First(&url).Error
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *URLRepository) IncrementClicks(ctx context.Context, shortCode string) error {
	return r.DB.WithContext(ctx).Model(&models.URL{}).
		Where("short_code = ?", shortCode).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", 1)).Error
}