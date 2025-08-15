package main

import (
	categoryDelivery "NeliQuiz/internal/features/category/delivery"
	categoryDomain "NeliQuiz/internal/features/category/domain"
	categoryRepository "NeliQuiz/internal/features/category/repository"
	categoryService "NeliQuiz/internal/features/category/usecase"
	questionDelivery "NeliQuiz/internal/features/question/delivery"
	questionDomain "NeliQuiz/internal/features/question/domain"
	questionRepository "NeliQuiz/internal/features/question/repository"
	questionService "NeliQuiz/internal/features/question/usecase"
	"NeliQuiz/internal/infrastructures/database/postgres"
	"NeliQuiz/internal/infrastructures/webserver"
	"NeliQuiz/internal/shared/config"
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
	cfg := config.New(".env")
	db, err := postgres.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := fx.New(
		fx.Provide(
			func() (*config.Config, *gorm.DB) {
				return cfg, db
			},
			fx.Annotate(categoryRepository.NewCategoryRepository, fx.As(new(categoryDomain.CategoryRepository))),
			fx.Annotate(questionRepository.NewQuestionRepository, fx.As(new(questionDomain.QuestionRepository))),
			fx.Annotate(categoryService.NewCategoryUseCase, fx.As(new(categoryDomain.CategoryUseCase))),
			fx.Annotate(questionService.NewQuestionUseCase, fx.As(new(questionDomain.QuestionUseCase))),
			categoryDelivery.NewCategoryHandler,
			questionDelivery.NewQuestionHandler,
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
