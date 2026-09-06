package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const (
	contextKeyRequestID contextKey = "request_id"
	contextKeyAccountID contextKey = "auth__account_id"
)

type responseCapture struct {
	http.ResponseWriter
	status int
}

func (r *responseCapture) WriteHeader(status int) {
	if r.status == 0 {
		r.status = status
		r.ResponseWriter.WriteHeader(status)
	}
}

func loggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			capture := &responseCapture{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(capture, r)

			logger.Info("http request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", capture.status,
				"duration", time.Since(start),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"referer", r.Referer(),
				"start_time", start.UTC().Format(time.RFC3339Nano),
				"end_time", time.Now().UTC().Format(time.RFC3339Nano),
				"request_id", r.Context().Value(contextKeyRequestID),
			)

		})

	}

}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), contextKeyRequestID, id)
		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")

		// Mocked -- not a real JWT token validation and parsing
		accountID, err := uuid.Parse(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyAccountID, accountID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
