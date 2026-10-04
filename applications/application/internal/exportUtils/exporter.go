package exportutils

import (
	"context"
	"io"
	"io/fs"
	"path/filepath"
)

type ExportFile struct {
	UUID     string
	Path     string
	FileInfo fs.FileInfo
	Name     string
}

type Exporter interface {
	Export(context.Context, io.Writer, []ExportFile) error
}

func (f ExportFile) ArchiveName() string {
	name := f.Name
	if name == "" {
		name = filepath.Base(f.Path)
	}
	return filepath.ToSlash(name)
}
