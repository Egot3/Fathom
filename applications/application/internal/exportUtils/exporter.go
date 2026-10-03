package exportutils

import (
	"context"
	"io"
	"io/fs"
	"time"
)

type ExportFile struct {
	UUID     string
	Path     string
	FileInfo fs.FileInfo
}

type customFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
}

func NewEmbeddedFile(name string, size int64) fs.FileInfo {
	return customFileInfo{
		name:    name,
		size:    size,
		mode:    0644,
		modTime: time.Now(),
	}
}

func (fi customFileInfo) Name() string       { return fi.name }
func (fi customFileInfo) Size() int64        { return fi.size }
func (fi customFileInfo) Mode() fs.FileMode  { return fi.mode }
func (fi customFileInfo) ModTime() time.Time { return fi.modTime }
func (fi customFileInfo) IsDir() bool        { return false }
func (fi customFileInfo) Sys() any           { return 1 }

type Exporter interface {
	Export(context.Context, io.Writer, []ExportFile) error
}
