package exportutlis

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"log/slog"

	"github.com/egot3/fathom/internal/logging"
)

type gzipExporter struct{}

func NewGzipExporter() Exporter {
	return &tarExporter{}
}

func (g *gzipExporter) Export(ctx context.Context, w io.Writer, files []ExportFile) error {
	logger := logging.LoggerFromContext(ctx).With(slog.String("layer", "exporter"))

	gzipWriter := gzip.NewWriter(w)
	tarWriter := tar.NewWriter(gzipWriter)

	for _, f := range files {
		if err := AddFileToTar(tarWriter, f.Path, f.FileInfo); err != nil {
			logger.Error("error writing to tar",
				slog.String("path", f.Path),
				slog.String("Error", err.Error()),
			)

			tarWriter.Close()
			return err
		}
	}

	if err := tarWriter.Close(); err != nil {
		logger.Error("error finalising tar", "error", err)
		return err
	}
	if err := gzipWriter.Close(); err != nil {
		logger.Error("error finalising gz", "error", err)
		return err
	}
	return nil
}
