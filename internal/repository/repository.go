package repository

import (
	"context"

	"github.com/MaksimovDenis/loadinator2000/internal/models"
)

type LoaderRepository interface {
	Create(ctx context.Context, fileName string, filePath string) (string, error)
	List(ctx context.Context, limit, offset int64) ([]models.FileInfo, error)
	Get(ctx context.Context, fileName string) (string, error)
}
