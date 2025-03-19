package loader

import (
	"context"
	"errors"
	"fmt"
	"os"
)

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
