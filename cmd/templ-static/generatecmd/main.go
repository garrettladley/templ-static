package generatecmd

import (
	"context"
	_ "embed"
	"log/slog"

	_ "net/http/pprof"
)

type Arguments struct {
	Path string
}

func Run(ctx context.Context, log *slog.Logger, args Arguments) (err error) {
	g, err := NewGenerate(log, args)
	if err != nil {
		return err
	}
	return g.Run(ctx)
}
