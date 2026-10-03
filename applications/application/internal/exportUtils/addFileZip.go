package exportutils

import (
	"archive/zip"
	"io"
	"os"
)

func AddFileToZip(zw *zip.Writer, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := zw.Create(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, f)
	return err
}
