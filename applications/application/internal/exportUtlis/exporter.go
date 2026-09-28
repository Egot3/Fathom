package exportutlis

import (
	"context"
	"io"
	"os"
)

type ExportFile struct {
	UUID     string
	Path     string
	FileInfo os.FileInfo // used by tar
}

type Exporter interface {
	Export(context.Context, io.Writer, []ExportFile) error
}
