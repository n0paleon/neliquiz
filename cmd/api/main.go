package main

import (
	"NeliQuiz/internal/config"
	"NeliQuiz/internal/delivery/http"
	"NeliQuiz/internal/domain"
	"NeliQuiz/internal/infrastructures/database/postgres"
	"NeliQuiz/internal/infrastructures/webserver"
	"NeliQuiz/internal/repository"
	"NeliQuiz/internal/usecase"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"net"
	"os"
)

// init logrus
func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	log.SetReportCaller(true)
}

func main() {
	cfg := config.New()
	db, err := postgres.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := fx.New(
		fx.Provide(
			func() (*config.Config, *gorm.DB) {
				return cfg, db
			},
			fx.Annotate(repository.NewPGQuestionRepository, fx.As(new(domain.QuestionRepository))),
			fx.Annotate(repository.NewPGCategoryRepository, fx.As(new(domain.CategoryRepository))),
			fx.Annotate(usecase.NewAdminQuestionUseCase, fx.As(new(domain.AdminQuestionUseCase))),
			fx.Annotate(usecase.NewUserQuestionUseCase, fx.As(new(domain.UserQuestionUseCase))),
			http.NewAdminQuestionHandler,
			http.NewUserQuestionHandler,
			webserver.NewRouter,
			webserver.NewServer,
		),
		fx.Invoke(func(r *webserver.Router, server *webserver.Server) error {
			r.RegisterRoutes()

			addr := net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort)
			return server.Listen(addr)
		}),
	)

	app.Run()
}
