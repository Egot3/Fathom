package acceptutils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func BestAccept(acceptHeader string, supported ...string) (string, error) {
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
			if strings.EqualFold(it.mime, sup) {
				return sup, nil
			}
		}
	}
	return "", fmt.Errorf("no supported type found")
}
