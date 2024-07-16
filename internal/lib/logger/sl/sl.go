package sl

import (
	"log/slog"
)

// github.com/mattn/go-sqlite3"
// Err - создание атрибута ошибки (имя параметра error, текстовое содержимое ошибки)
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
