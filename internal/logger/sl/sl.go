package sl

import (
	"log/slog"

	_ "github.com/lib/pq"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "err",
		Value: slog.StringValue(err.Error()),
	}
}
