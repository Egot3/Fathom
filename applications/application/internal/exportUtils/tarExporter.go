package exportutils

import (
	"archive/tar"
	"context"
	"io"
	"log/slog"

	"github.com/egot3/fathom/internal/logging"
)

type tarExporter struct{}

func NewTarExporter() Exporter {
	return &tarExporter{}
}

func (z *tarExporter) Export(ctx context.Context, w io.Writer, files []ExportFile) error {
	logger := logging.LoggerFromContext(ctx).With(slog.String("layer", "exporter"))

	tarWriter := tar.NewWriter(w)

	for _, f := range files {
		if err := AddFileToTar(tarWriter, f); err != nil {
			logger.Error("error writing to tar",
				slog.String("path", f.Path),
				slog.String("Error", err.Error()),
			)

			tarWriter.Close()
			return err
		}
	}

	err := tarWriter.Close()
	if err != nil {
		logger.Error("error finalising tar", slog.String("Error", err.Error()))
	}
	return err
}
