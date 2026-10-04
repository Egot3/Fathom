package httputils

import (
	"sort"
	"strconv"
	"strings"

	"github.com/egot3/fathom/internal/carefulness"
)

type EIMimes = string

const (
	Tar  EIMimes = "application/tar"
	GZip EIMimes = "application/gzip"
	Zip  EIMimes = "application/zip"
	Yaml EIMimes = "application/yaml"
)

var AvailableArchiveMimes = []EIMimes{Tar, GZip, Zip, Yaml}

func BestAccept[T ~string](acceptHeader string, supported ...T) (T, carefulness.JSONErrorable) {
	if acceptHeader == "" {
		if len(supported) > 0 {
			return supported[0], nil
		}
		return "", nil
	}

	type media struct {
		mime    string
		quality float32
	}
	var items []media

	for part := range strings.SplitSeq(acceptHeader, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		segments := strings.Split(part, ";")
		mime := strings.TrimSpace(segments[0])
		q := float32(1.0)
		for _, seg := range segments[1:] {
			seg = strings.TrimSpace(seg)
			if strings.HasPrefix(seg, "q=") {
				val, err := strconv.ParseFloat(seg[2:], 32)
				if err == nil {
					q = float32(val)
				}
			}
		}
		if q > 0 {
			items = append(items, media{mime, q})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].quality > items[j].quality
	})

	for _, it := range items {
		for _, sup := range supported {
			if strings.EqualFold(it.mime, string(sup)) {
				return sup, nil
			}
		}
	}
	return "", carefulness.ErrUnnacaptable
}
