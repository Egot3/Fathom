package exportutils

import (
	"archive/zip"
	"io"
	"os"
)

func AddFileToZip(zw *zip.Writer, f ExportFile) error {
	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	w, err := zw.Create(f.Path)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, file)
	return err
}
