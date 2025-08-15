package webserver

import (
	"NeliQuiz/internal/shared/apihelper"
	"NeliQuiz/internal/shared/config"
	"NeliQuiz/internal/shared/errorx"
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
		ErrorHandler:      ErrorHandler,
		BodyLimit:         4 * 1024 * 1024, // 4 MB
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ProxyHeader:       fiber.HeaderXForwardedFor,
		CaseSensitive:     true,
		StrictRouting:     true,
		AppName:           "NeliQuiz",
		JSONEncoder:       sonic.Marshal,
		JSONDecoder:       sonic.Unmarshal,
		EnablePrintRoutes: true,
	})

	if cfg.ValidateApiGateway {
		app.Use(ApiGatewayAuthMiddleware(cfg.ApiGatewayToken))
	}
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

func ApiGatewayAuthMiddleware(secretToken string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Api-Gateway-Token", "")
		if token == "" || token != secretToken {
			// block connection
			_ = c.Context().Conn().Close()
			return nil
		}

		// continue to next route
		return c.Next()
	}
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

		return c.Status(appErr.Code()).JSON(apihelper.APIResponse[any]{
			Status:  "error",
			Data:    nil,
			Message: appErr.Message(),
		})
	}

	// Unhandled error
	logrus.WithFields(map[string]interface{}{
		"error":  err.Error(),
		"path":   c.Path(),
		"method": c.Method(),
	}).Error("unhandled error")

	return c.Status(fiber.StatusInternalServerError).JSON(apihelper.APIResponse[any]{
		Status:  "error",
		Data:    nil,
		Message: "internal server error",
	})
}
