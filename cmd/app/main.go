package main

import (
	categoryDelivery "NeliQuiz/internal/features/category/delivery"
	categoryDomain "NeliQuiz/internal/features/category/domain"
	categoryRepository "NeliQuiz/internal/features/category/repository"
	categoryService "NeliQuiz/internal/features/category/usecase"
	"NeliQuiz/internal/features/game/delivery/modes/broadcast"
	"NeliQuiz/internal/features/game/delivery/modes/solo"
	questionDelivery "NeliQuiz/internal/features/question/delivery"
	questionDomain "NeliQuiz/internal/features/question/domain"
	questionRepository "NeliQuiz/internal/features/question/repository"
	questionService "NeliQuiz/internal/features/question/usecase"
	"NeliQuiz/internal/infrastructures/database/postgres"
	"NeliQuiz/internal/infrastructures/gameserver"
	"NeliQuiz/internal/infrastructures/webserver"
	"NeliQuiz/internal/infrastructures/workerpool"
	"NeliQuiz/internal/shared/config"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// init logrus
func init() {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
}

func main() {
	cfg := config.New(".env")
	db, err := postgres.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	wp := workerpool.InitPool(10000)

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
			gameserver.NewGameServer,
			solo.NewSoloRoom,
			broadcast.NewBroadcastRoom,
		),
		fx.Invoke(func(r *webserver.Router, web *webserver.Server) error {
			r.RegisterRoutes()
			return nil
		}),
		fx.Invoke(func(solo *solo.Room, gameServer *gameserver.Server, broadcastRoom *broadcast.Room) error {
			gameServer.RegisterComponent(solo,
				component.WithName("solo"),
			)
			gameServer.RegisterComponent(broadcastRoom,
				component.WithName("broadcast"),
			)
			return nil
		}),
		fx.Invoke(func(lc fx.Lifecycle, webServer *webserver.Server, gameServer *gameserver.Server) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Info("🚀 Starting NeliQuiz application...")

					// Start web server
					_ = wp.Submit(func() {
						log.Info("🌐 Starting web server...")
						addr := net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort)
						if err := webServer.Listen(addr); err != nil {
							log.Fatalf("Web server startup failed: %v", err)
						}
					})

					// Start game server
					_ = wp.Submit(func() {
						log.Info("🎮 Starting game server...")
						gameServer.Start(ctx)
					})

					log.Info("✅ All servers started successfully")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("🛑 Shutting down NeliQuiz application...")

					// Stop web server gracefully
					log.Info("🌐 Stopping web server...")
					if err := webServer.Shutdown(ctx); err != nil {
						log.Errorf("Web server shutdown error: %v", err)
					}

					// Close database connection
					log.Info("💾 Closing database connection...")
					if sqlDB, err := db.DB(); err == nil {
						_ = sqlDB.Close()
					}

					log.Info("✅ Application shutdown completed")
					return nil
				},
			})
		}),
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// wait for quit signal
	<-quit
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := app.Stop(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
