package loader

import (
	"context"

	"github.com/MaksimovDenis/loadinator2000/internal/models"
)

func (srv *serv) List(ctx context.Context, limit, offset int64) ([]models.FileInfo, error) {
	select {
	case srv.listLimiter <- struct{}{}:
		defer func() { <-srv.listLimiter }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if limit <= 0 {
		limit = 100
	}

	if offset <= 0 {
		offset = 0
	}

	fileInfo, err := srv.loaderRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}
