package loader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var invalidCharsRegex = regexp.MustCompile(`[\"'<>!@#$%^&*()=+\[\]{}|\\/]`)

func (srv *serv) Create(ctx context.Context, fileName string, filePath string, data []byte) (string, error) {
	select {
	case srv.downloadLimiter <- struct{}{}:
		defer func() { <-srv.downloadLimiter }()
	case <-ctx.Done():
		return "", ctx.Err()
	}

	if err := validateData(fileName, filePath, data); err != nil {
		return "", err
	}

	var fileNameResp string

	err := srv.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		_, errTx = srv.loaderRepository.Get(ctx, fileName)
		if errTx == nil {
			return fmt.Errorf("файл с таким именем уже существует")
		}

		fullPath, errTx := storeFile(fileName, filePath, data)
		if errTx != nil {
			return errTx
		}

		fileNameResp, errTx = srv.loaderRepository.Create(ctx, fileName, fullPath)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return fileNameResp, nil
}

func validateData(fileName, filePath string, data []byte) error {
	switch {
	case fileName == "":
		return errors.New("заполните поле имя файла")
	case filePath == "":
		return errors.New("заполните поле путь файла")
	case len(data) == 0:
		return errors.New("бинарное содержимое файла пустое")
	case invalidCharsRegex.MatchString(fileName):
		return errors.New("имя файла содержит недопустимые символы")
	default:
		return nil
	}
}

func storeFile(fileName, filePath string, data []byte) (string, error) {
	filePath = strings.TrimRight(filePath, "/")

	if strings.HasPrefix(filePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("не удалось получить домашнюю директорию: %v", err)
		}

		filePath = filepath.Join(homeDir, filePath[1:])
	}

	fullPath := filepath.Join(filePath, fileName)

	dir := filepath.Dir(fullPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("не удалось создать директорию %s: %v", dir, err)
		}
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл %s: %v", fullPath, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return "", fmt.Errorf("не удалось записать данные в файл %s: %v", fullPath, err)
	}

	return fullPath, nil
}
