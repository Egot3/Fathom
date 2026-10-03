package exportutils

import (
	"archive/tar"
	"io"
	"os"
)

func AddFileToTar(tw *tar.Writer, f ExportFile) error {
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := tar.FileInfoHeader(f.FileInfo, "")
	if err != nil {
		return err
	}
	header.Name = f.ArchiveName()

	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	_, err = io.Copy(tw, file)
	return err
}
