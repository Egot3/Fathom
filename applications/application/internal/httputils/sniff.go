package httputils

import (
	"bytes"
	"io"
	"mime/multipart"
)

var magicBytesTar = [][]byte{
	[]byte("ustar\x00tar\x00"),
	[]byte("ustar\x00"),
	[]byte("ustar \x00"),
}

// https://media1.giphy.com/media/v1.Y2lkPTc5MGI3NjExcWl2aW9wemw0OGRkbGdiZnh6MXdtZnZlbTBkbXVhZXZ1Y2Zuanc0dSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/hpqWzmAXaaBl1heu3i/giphy.gif
func SniffArchive(f multipart.File) (EIMimes, bool) {
	var head [8]byte
	n, _ := io.ReadFull(f, head[:])
	defer f.Seek(0, io.SeekStart)

	if n < 4 {
		return "", false
	}

	switch {
	case head[0] == 'P' && head[1] == 'K' && head[2] == 3 && head[3] == 4:
		return Zip, true
	case head[0] == 0x1f && head[1] == 0x8b:
		return GZip, true
	default:
		_, err := f.Seek(257, io.SeekStart)
		if err != nil {
			return "", false
		}

		n, _ = io.ReadFull(f, head[:])
		if n == 0 {
			return "", false
		}

		for _, magic := range magicBytesTar {
			if bytes.HasPrefix(head[:n], magic) {
				return Tar, true
			}
		}
	}

	return "", false
}
