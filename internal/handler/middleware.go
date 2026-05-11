package handler

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/yloveya1/metricsalert/internal/logger"
	"go.uber.org/zap"
)

const (
	hashHeader    = "HashSHA256"
	contentHeader = "Content-Type"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			ow.Header().Add("Content-Encoding", "gzip")
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}

func (h *Handler) WithLogging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				logger.ServerLog.Info(
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

func (h *Handler) HashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Get(contentHeader)) == 0 || len(h.key) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		hashHeaderValue := r.Header.Get(hashHeader)
		if hashHeaderValue == "" {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(body))

		hasher := hmac.New(sha256.New, []byte(h.key))
		hasher.Write(body)
		computedMAC := hasher.Sum(nil)

		decodedMAC, err := hex.DecodeString(hashHeaderValue)
		if err != nil {
			http.Error(w, "Invalid hash encoding", http.StatusBadRequest)
			return
		}

		if !hmac.Equal(computedMAC, decodedMAC) {
			http.Error(w, "Invalid hash", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
