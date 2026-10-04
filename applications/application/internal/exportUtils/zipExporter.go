package exportutils

import (
	"archive/zip"
	"context"
	"io"
	"log/slog"

	"github.com/egot3/fathom/internal/logging"
)

type zipExporter struct{}

func NewZipExporter() Exporter {
	return &zipExporter{}
}

func (z *zipExporter) Export(ctx context.Context, w io.Writer, files []ExportFile) error {
	logger := logging.LoggerFromContext(ctx).With(slog.String("layer", "exporter"))

	zipWriter := zip.NewWriter(w)

	for _, f := range files {
		if err := AddFileToZip(zipWriter, f); err != nil {
			logger.Error("error writing to zip",
				slog.String("path", f.Path),
				slog.String("Error", err.Error()),
			)
			zipWriter.Close()
			return err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		logger.Error("error finalising zip", slog.String("Error", err.Error()))
	}
	return err
}
