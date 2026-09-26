package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/egot3/fathom/version"
)

type Status string

const (
	StatusUp   Status = "up"
	StatusDown Status = "down"
)

type Check struct {
	Name string
	Fn   func(ctx context.Context) error
}

type component struct {
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Response struct {
	Status     Status               `json:"status"`
	Version    version.Version      `json:"version,omitempty"`
	Components map[string]component `json:"components,omitempty"`
}

func Handler(version version.Version, timeout time.Duration, checks ...Check) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		resp := Response{Status: StatusUp, Version: version, Components: make(map[string]component, len(checks))}

		var mu sync.Mutex
		var wg sync.WaitGroup
		for _, c := range checks {
			wg.Add(1)
			go func(c Check) {
				defer wg.Done()
				err := c.Fn(ctx)

				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					resp.Status = StatusDown
					resp.Components[c.Name] = component{Status: StatusDown, Error: err.Error()}
					return
				}
				resp.Components[c.Name] = component{Status: StatusUp}
			}(c)
		}
		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		if resp.Status == StatusDown {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		json.NewEncoder(w).Encode(resp)
	}
}
