package archiveutlis

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func AddFileToTar(tw *tar.Writer, filePath string, fi os.FileInfo) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	header, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filePath)

	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	_, err = io.Copy(tw, f)
	return err
}
