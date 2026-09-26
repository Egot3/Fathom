package middlewares

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/egot3/fathom/internal/logging"
)

func Recoverer(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					panic(rvr)
				}

				logger := logging.LoggerFromContext(r.Context())
				logger.Error("panic!", slog.String("stack", string(debug.Stack())), slog.Any("recover", rvr))

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
