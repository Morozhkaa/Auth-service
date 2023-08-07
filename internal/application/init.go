package application

import (
	"auth/internal/adapters/db"
	"auth/internal/adapters/http"
	"auth/internal/adapters/memory"
	"auth/internal/config"
	"auth/internal/domain/usecases"
	"auth/pkg/infra/logger"
	"context"
)

type App struct {
	logger        logger.Logger
	shutdownFuncs []func(ctx context.Context) error
}

func New(logger logger.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (app *App) Start() error {
	cfg, _ := config.GetConfig()
	var auth *usecases.Auth

	if cfg.DB_used {
		userStorage, err := db.New(context.Background(), cfg.DB_url)
		if err != nil {
			app.logger.Sugar().Fatalf("create user storage failed: %s", err.Error())
		}
		auth, err = usecases.New(userStorage)
		if err != nil {
			app.logger.Sugar().Fatalf("create buissness logic failed: %s", err.Error())
		}
	} else {
		userStorage, err := memory.New()
		if err != nil {
			app.logger.Sugar().Fatalf("create user storage failed: %s", err.Error())
		}
		auth, err = usecases.New(userStorage)
		if err != nil {
			app.logger.Sugar().Fatalf("create buissness logic failed: %s", err.Error())
		}
	}

	s, err := http.New(auth, app.logger)
	if err != nil {
		app.logger.Sugar().Fatalf("server not started %s", err.Error())
	}
	app.shutdownFuncs = append(app.shutdownFuncs, s.Stop)
	err = s.Start()
	if err != nil {
		app.logger.Sugar().Fatalf("server not started: %s", err.Error())
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	var err error
	for i := len(a.shutdownFuncs) - 1; i >= 0; i-- {
		err = a.shutdownFuncs[i](ctx)
		if err != nil {
			a.logger.Sugar().Error(err)
		}
	}
	a.logger.Info("app stopped")
	return nil
}
