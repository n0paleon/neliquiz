package webserver

import (
	"NeliQuiz/internal/config"
	"NeliQuiz/internal/errorx"
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"time"
)

type Server struct {
	app *fiber.App
}

func NewServer(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		// Error handler global
		ErrorHandler: ErrorHandler,
		// Maksimum ukuran body (hindari payload membabi buta)
		BodyLimit: 4 * 1024 * 1024, // 4 MB
		// Maksimum waktu baca header (anti slowloris)
		ReadTimeout: 10 * time.Second,
		// Maksimum waktu tulis response
		WriteTimeout: 10 * time.Second,
		// Proxy trusted (kalau pakai reverse proxy, load balancer)
		// misal TrueClientIP / X-Forwarded-For
		ProxyHeader: fiber.HeaderXForwardedFor,
		// Case sensitive URL: true = /users ≠ /Users
		CaseSensitive: true,
		// Strict routing: true = /users ≠ /users/
		StrictRouting: true,
		AppName:       "NeliQuiz",
		// JSON Encoder
		JSONEncoder: sonic.Marshal,
		// JSON Decoder
		JSONDecoder: sonic.Unmarshal,
		// Enable print routes?
		EnablePrintRoutes: true,
	})

	app.Use(recover2.New())

	return &Server{app: app}
}

func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func (s *Server) ShutdownWithContext(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	var appErr errorx.AppError
	if errors.As(err, &appErr) {
		if appErr.Code() >= 500 {
			logrus.WithFields(map[string]interface{}{
				"cause":  appErr.Cause(),
				"path":   c.Path(),
				"method": c.Method(),
			}).Error(appErr.Message())
		}

		return c.Status(appErr.Code()).JSON(fiber.Map{
			"status": "error",
			"data":   "",
			"error":  appErr.Message(),
		})
	}

	// Unhandled error
	logrus.WithFields(map[string]interface{}{
		"error":  err.Error(),
		"path":   c.Path(),
		"method": c.Method(),
	}).Error("unhandled error")

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status": "error",
		"data":   "",
		"error":  "internal server error",
	})
}
