package hashutils

import (
	"io"
	"os"

	"github.com/zeebo/xxh3"
)

func HashFile(file *os.File) (uint64, error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	hasher := xxh3.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return 0, err
	}

	return hasher.Sum64(), nil
}
