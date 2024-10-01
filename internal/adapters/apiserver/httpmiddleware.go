package apiserver

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"kl-webhook-adapter/internal/core/app"
)

const httpLogFormat = `"[END] %s %s %s" from %s`
const xRequestId = "x-request-id"

// zapLoggerMiddleware provides middleware for adding zap logger into the context of request handler
func zapLoggerMiddleware(logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqId := correlationId(r)
			w.Header().Set(xRequestId, reqId)
			l := logger.With(zap.String("requestId", reqId))
			ctx := app.ContextWithLogger(r.Context(), l)

			tbegin := time.Now()
			defer func() {
				l.With(zap.String("duration", time.Since(tbegin).String())).Infof(httpLogFormat,
					r.Method,
					r.URL.Path,
					r.Proto,
					r.RemoteAddr,
				)
			}()
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func correlationId(r *http.Request) string {
	reqId := r.Header.Get(xRequestId)
	if reqId == "" {
		reqId = uuid.NewString()
	}
	return reqId
}

func recoveryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					app.Logger(r.Context()).Errorln("Recovered from panic:", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
