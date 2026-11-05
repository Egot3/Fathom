package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

func loadIndex(fsys fs.FS) ([]byte, error) {
	raw, err := fs.ReadFile(fsys, "index.html")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	cfg := map[string]string{}
	if v := os.Getenv("PUBLIC_BASE_URL"); v != "" {
		cfg["publicBaseUrl"] = strings.TrimRight(v, "/")
	}
	j, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	tag := "<script>window.__CONFIG__=" + string(j) + "</script></head>"
	return bytes.Replace(raw, []byte("</head>"), []byte(tag), 1), nil
}

func SPAHandler(fsys fs.FS) (http.Handler, error) {
	index, err := loadIndex(fsys)
	if err != nil {
		return nil, err
	}
	files := http.FileServer(http.FS(fsys))

	serveIndex := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(index)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")

		if p == "" || p == "index.html" {
			serveIndex(w, r)
			return
		}

		if f, err := fsys.Open(p); err == nil {
			st, statErr := f.Stat()
			f.Close()
			if statErr == nil && !st.IsDir() {
				if strings.HasPrefix(p, "assets/") {
					w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
				} else {
					w.Header().Set("Cache-Control", "no-cache")
				}
				files.ServeHTTP(w, r)
				return
			}
		}

		if path.Ext(p) != "" {
			http.NotFound(w, r)
			return
		}

		serveIndex(w, r)
	}), nil
}
