package httpapi

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

type Config struct {
	Addr string
	User string
	Pass string
}

type Controller struct {
	logger *slog.Logger
	cfg    Config
}

func New(logger *slog.Logger, cfg Config) Controller {
	return Controller{
		logger: logger,
		cfg:    cfg,
	}
}

func (cnt Controller) Serve(ctx context.Context) error {
	e := echo.New()
	e.Validator = vldr{validator: validator.New()}

	e.Use(
		slogecho.NewWithConfig(cnt.logger, slogecho.Config{
			DefaultLevel:     slog.LevelDebug,
			ClientErrorLevel: slog.LevelWarn,
			ServerErrorLevel: slog.LevelError,

			WithUserAgent:      true,
			WithRequestID:      true,
			WithRequestBody:    true,
			WithRequestHeader:  true,
			WithResponseBody:   true,
			WithResponseHeader: true,
			WithClientIP:       true,
		}),
		middleware.CORS(),
		middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
			Skipper: func(c echo.Context) bool {
				if c.Path() == "/metrics" {
					return true
				}

				if cnt.cfg.User == "" || cnt.cfg.Pass == "" {
					return true
				}

				return false
			},
			Realm: "easyjet",
			Validator: func(s1, s2 string, ctx echo.Context) (bool, error) {
				if cnt.cfg.User == s1 && cnt.cfg.Pass == s2 {
					return true, nil
				}

				return false, nil
			},
		}),
	)

	go func() {
		<-ctx.Done()
		err := e.Shutdown(context.Background())
		if err != nil {
			cnt.logger.Error("shutdown http", slog.String("error", err.Error()))
		}
	}()

	err := e.Start(cnt.cfg.Addr)
	if err != nil {
		return err
	}

	return nil
}
