package loader

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// Get загружает файл по его имени.
//
// 1. Проверяет лимит на количество одновременных операций (listLimiter).
// 2. Проверяет, не был ли контекст отменен перед началом выполнения.
// 3. Валидирует имя файла.
// 4. Ищет путь к файлу в базе данных по имени файла.
// 5. Загружает содержимое файла с диска по найденному пути.
// 6. Возвращает содержимое файла или ошибку в случае неудачи.
func (srv *serv) Get(ctx context.Context, fileName string) ([]byte, error) {
	select {
	case srv.listLimiter <- struct{}{}:
		defer func() { <-srv.listLimiter }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	err := validateFileName(fileName)
	if err != nil {
		return nil, err
	}

	filePath, err := srv.loaderRepository.Get(ctx, fileName)
	if err != nil {
		return nil, err
	}

	data, err := loadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func validateFileName(fileName string) error {
	switch {
	case fileName == "":
		return errors.New("заполните поле имя файла")
	case invalidCharsRegex.MatchString(fileName):
		return errors.New("имя файла содержит недопустимые символы")
	default:
		return nil
	}
}

func loadFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %s %w", filepath, err)
	}
	defer file.Close()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл %s %w", filepath, err)
	}

	return data, nil
}
