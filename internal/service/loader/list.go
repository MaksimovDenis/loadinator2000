package loader

import (
	"context"

	"github.com/MaksimovDenis/loadinator2000/internal/models"
)

// List возвращает список загруженных файлов с учетом пагинации.
//
// 1. Проверяет лимит на количество одновременных операций (listLimiter).
// 2. Проверяет, не был ли контекст отменен перед началом выполнения.
// 3. Устанавливает значение `limit` по умолчанию (100), если оно некорректно.
// 4. Устанавливает значение `offset` по умолчанию (0), если оно некорректно.
// 5. Запрашивает список файлов из хранилища с учетом `limit` и `offset`.
// 6. Возвращает список файлов или ошибку в случае неудачи.
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
