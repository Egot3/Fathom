package hashutils

import (
	"crypto/md5"
	"io"
	"os"
)

func HashFileMD5(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hasher := md5.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	hashBytes := hasher.Sum(nil)

	return hashBytes, nil
}
