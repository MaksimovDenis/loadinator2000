package models

import (
	"database/sql"
	"time"
)

type FileInfo struct {
	FileName  string       `json:"filename" db:"filename"`
	FilePath  string       `json:"file_path" db:"file_path"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"`
}
