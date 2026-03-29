package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/yloveya1/metricsalert/internal/logger"
	"go.uber.org/zap"
)

func (h *Handler) WithLogging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				logger.Log.Info(
					"got incoming HTTP request",
					zap.String("uri", r.RequestURI),
					zap.String("method", r.Method),
					zap.String("duration", time.Since(t1).String()),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
