package loader

import (
	"context"

	db "github.com/MaksimovDenis/loadinator2000/internal/client"
	"github.com/MaksimovDenis/loadinator2000/internal/models"
	"github.com/MaksimovDenis/loadinator2000/internal/repository"

	"github.com/Masterminds/squirrel"
	"github.com/rs/zerolog"
)

const (
	tableName = "files"
)

type repo struct {
	db  db.Client
	log zerolog.Logger
}

func NewRepository(db db.Client, log zerolog.Logger) repository.LoaderRepository {
	return &repo{
		db:  db,
		log: log,
	}
}

func (rep *repo) Create(ctx context.Context, fileName string, filePath string) (string, error) {
	builder := squirrel.Insert(tableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns("filename", "file_path").
		Values(fileName, filePath).
		Suffix("RETURNING filename")

	query, args, err := builder.ToSql()
	if err != nil {
		rep.log.Error().Err(err).Msg("create file: failed to build SQL query")
		return "", err
	}

	queryStruct := db.Query{
		Name:     "loader_repository.AddTransaction",
		QueryRow: query,
	}

	var name string
	err = rep.db.DB().QueryRowContext(ctx, queryStruct, args...).Scan(&name)
	if err != nil {
		rep.log.Error().Err(err).Msg("failed to create file")
		return "", err
	}

	return name, nil
}

func (rep *repo) List(ctx context.Context, limit, offset int64) ([]models.FileInfo, error) {
	builder := squirrel.Select("filename, file_path, created_at, updated_at").
		PlaceholderFormat(squirrel.Dollar).
		From(tableName).
		Limit(uint64(limit)).
		Offset(uint64(offset))

	query, args, err := builder.ToSql()
	if err != nil {
		rep.log.Error().Err(err).Msg("get files list: failed to build SQL query")
		return nil, err
	}

	queryStruct := db.Query{
		Name:     "loader_repository.get files list",
		QueryRow: query,
	}

	var fileInfo []models.FileInfo

	err = rep.db.DB().ScanAllContext(ctx, &fileInfo, queryStruct, args...)
	if err != nil {
		rep.log.Error().Err(err).Msg("get files list: failed to scan rows")
		return nil, err
	}

	return fileInfo, nil
}

func (rep *repo) Get(ctx context.Context, fileName string) (string, error) {
	builder := squirrel.Select("file_path").
		PlaceholderFormat(squirrel.Dollar).
		From(tableName).
		Where(squirrel.Eq{"filename": fileName}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		rep.log.Error().Err(err).Msg("get file: failed to build SQL query")
		return "", err
	}

	queryStruct := db.Query{
		Name:     "loader_repository.get file",
		QueryRow: query,
	}

	var filePath string

	err = rep.db.DB().QueryRowContext(ctx, queryStruct, args...).Scan(&filePath)
	if err != nil {
		rep.log.Error().Err(err).Msg("get file: failed to scan row")
		return "", err
	}

	return filePath, nil
}
