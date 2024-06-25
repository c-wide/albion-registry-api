package server

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func registerMiddleware(e *echo.Echo, logger zerolog.Logger) {
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.Error().Err(err).Msg("Recovered from panic in HTTP handler")
			return err
		},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			if strings.Contains(c.Request().URL.Path, "swagger") {
				return true
			}
			return false
		},
		Level: 5,
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:      true,
		LogRemoteIP:     true,
		LogHost:         true,
		LogMethod:       true,
		LogURI:          true,
		LogUserAgent:    true,
		LogStatus:       true,
		LogError:        true,
		LogResponseSize: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			var event *zerolog.Event

			if v.Error != nil {
				event = logger.Error().Err(v.Error)
			} else {
				event = logger.Info()
			}

			event.
				Dur("latency", v.Latency).
				Str("remoteIP", v.RemoteIP).
				Str("host", v.Host).
				Str("method", v.Method).
				Str("uri", v.URI).
				Str("userAgent", v.UserAgent).
				Int("status", v.Status).
				Int64("responseSize", v.ResponseSize)

			event.Send()

			return nil
		},
	}))

}
