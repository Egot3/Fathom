package importutils

import (
	"archive/tar"
	"context"
	"log/slog"
	"mime/multipart"

	"github.com/egot3/fathom/internal/httputils"
	"github.com/egot3/fathom/internal/logging"
)

type tarImporter struct{}

func NewTarImporter() Importer {
	return tarImporter{}
}

func (_ tarImporter) Import(ctx context.Context, r multipart.File, size int64, stageDir string, turnToAbs func(string) (string, error)) ([]StagedQuiz, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.String("strategy", httputils.Tar))
	ctx = logging.WithLogger(ctx, logger)

	return importTar(tar.NewReader(r), ctx, stageDir, turnToAbs)
}
