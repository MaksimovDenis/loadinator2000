package service

import (
	"context"

	"github.com/MaksimovDenis/loadinator2000/internal/models"
)

type LoaderService interface {
	Create(ctx context.Context, fileName string, filePath string, data []byte) (string, error)
	List(ctx context.Context, limit, offset int64) ([]models.FileInfo, error)
	Get(ctx context.Context, fileName string) ([]byte, error)
}
