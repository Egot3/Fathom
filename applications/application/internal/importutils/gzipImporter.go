package importutils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"log/slog"
	"mime/multipart"

	"github.com/egot3/fathom/internal/httputils"
	"github.com/egot3/fathom/internal/logging"
)

type gzipImporter struct{}

func NewGzipImporter() Importer {
	return gzipImporter{}
}

func (z gzipImporter) Import(ctx context.Context, r multipart.File, size int64, stageDir string, turnToAbs func(string) (string, error)) ([]StagedQuiz, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.String("strategy", httputils.GZip))
	ctx = logging.WithLogger(ctx, logger)

	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		logger.Error("couldn't create gzip reader",
			slog.String("Error", err.Error()),
		)
		return nil, err
	}
	tarReader := tar.NewReader(gzipReader)

	return importTar(tarReader, ctx, stageDir, turnToAbs)
}
