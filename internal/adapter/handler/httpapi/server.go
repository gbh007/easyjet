package httpapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gbh007/easyjet/internal/adapter/handler/httpapi/ogenapi"
	"github.com/gbh007/easyjet/internal/core/port"
	"github.com/gbh007/easyjet/pkg/metrics"
	"github.com/ogen-go/ogen/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Addr            string
	User            string
	Pass            string
	StaticFilesPath string
}

type Controller struct {
	logger  *slog.Logger
	cfg     Config
	service port.Service
}

func New(logger *slog.Logger, cfg Config, service port.Service) Controller {
	return Controller{
		logger:  logger,
		cfg:     cfg,
		service: service,
	}
}

func (cnt Controller) Serve(ctx context.Context) error {
	handler := NewHandler(cnt.service)

	server, err := ogenapi.NewServer(handler,
		ogenapi.WithMiddleware(middlewareFunc(cnt.logger)),
	)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	if cnt.cfg.StaticFilesPath != "" {
		mux.Handle("/", http.FileServer(http.Dir(cnt.cfg.StaticFilesPath)))
	}

	mux.Handle("/metrics/", promhttp.HandlerFor(metrics.DefaultRegistry, promhttp.HandlerOpts{}))
	mux.Handle("/api/", server)

	srv := &http.Server{
		Addr:    cnt.cfg.Addr,
		Handler: cnt.authMiddleware(mux),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		//nolint:contextcheck // shutdown requires new context after original is cancelled
		if err := srv.Shutdown(shutdownCtx); err != nil {
			cnt.logger.Error("shutdown http", slog.String("error", err.Error()))
		}
	}()

	err = srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func middlewareFunc(logger *slog.Logger) middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		// Log request
		logger.Debug("HTTP request",
			slog.String("method", req.Raw.Method),
			slog.String("path", req.Raw.URL.Path),
			slog.String("remote_addr", req.Raw.RemoteAddr),
			slog.String("operation", req.OperationName),
		)

		resp, err := next(req)

		if err != nil {
			logger.Error("HTTP response error",
				slog.String("method", req.Raw.Method),
				slog.String("path", req.Raw.URL.Path),
				slog.String("operation", req.OperationName),
				slog.String("error", err.Error()),
			)
		} else {
			logger.Debug("HTTP response",
				slog.String("method", req.Raw.Method),
				slog.String("path", req.Raw.URL.Path),
				slog.String("operation", req.OperationName),
			)
		}

		return resp, err
	}
}

func (cnt Controller) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		if cnt.cfg.User == "" || cnt.cfg.Pass == "" {
			next.ServeHTTP(w, r)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok || cnt.cfg.User != username || cnt.cfg.Pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="easyjet"`)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte(`{"error":"Unauthorized"}`))
			if err != nil {
				cnt.logger.Debug("write response error", slog.String("error", err.Error()))
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
