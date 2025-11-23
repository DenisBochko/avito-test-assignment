package app

import (
	"fmt"

	"go.uber.org/zap"

	"avito-test-assignment/internal/api/http/handler"
	"avito-test-assignment/internal/api/http/route"
	"avito-test-assignment/internal/config"
	"avito-test-assignment/internal/repository"
	"avito-test-assignment/internal/service"
	"avito-test-assignment/pkg/postgres"
	"avito-test-assignment/pkg/server"
)

type App struct {
	l          *zap.Logger
	cfg        *config.Config
	db         postgres.Postgres
	httpServer server.HTTPServer
}

type Repository struct {
	UserRepo *repository.UserRepository
	TeamRepo *repository.TeamRepository
}

type Service struct {
	TeamSvc *service.TeamService
	UserSvc *service.UserService
}

type Handler struct {
	TeamHdl *handler.TeamHandler
	UserHdl *handler.UserHandler
}

func New(l *zap.Logger, cfg *config.Config) (*App, error) {
	db, err := initDB(l, &cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	repo := initRepository(l, db)

	svc := initService(l, repo)

	hdl := initHandler(l, svc)

	httpServer := initHTTPServer(l, cfg, hdl)

	return &App{
		cfg:        cfg,
		l:          l,
		db:         db,
		httpServer: httpServer,
	}, nil
}

func MustNew(l *zap.Logger, cfg *config.Config) *App {
	app, err := New(l, cfg)
	if err != nil {
		panic(err)
	}

	return app
}

func (a *App) Run() error {
	errs := make(chan error, 1)
	defer close(errs)

	go func() {
		if err := a.httpServer.Run(); err != nil {
			errs <- err
		}
	}()

	if err := <-errs; err != nil {
		return err
	}

	return nil
}

func (a *App) Shutdown() error {
	a.db.Close()
	a.l.Debug("Database closed")

	if srvErr := a.httpServer.Shutdown(); srvErr != nil {
		return srvErr
	}

	a.l.Debug("HTTP server shutdown")

	return nil
}

func initDB(l *zap.Logger, cfg *config.Database) (postgres.Postgres, error) {
	postgresCfg := &postgres.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		Name:     cfg.Name,
		SSLMode:  cfg.SSLMode,
		MaxConns: cfg.MaxConns,
		MinConns: cfg.MinConns,
		Migration: postgres.Migration{
			Path:      cfg.Migration.Path,
			AutoApply: cfg.Migration.AutoApply,
		},
	}

	db, err := postgres.New(postgresCfg)
	if err != nil {
		return nil, err
	}

	l.Debug("Postgres initialized")

	return db, nil
}

func initRepository(l *zap.Logger, db postgres.Postgres) *Repository {
	userRepo := repository.NewUserRepository(db.Pool())

	l.Debug("User repository initialized")

	teamRepo := repository.NewTeamRepository(db.Pool())

	l.Debug("Team repository initialized")

	return &Repository{
		UserRepo: userRepo,
		TeamRepo: teamRepo,
	}
}

func initService(l *zap.Logger, repo *Repository) *Service {
	teamSvc := service.NewTeamService(repo.TeamRepo, repo.UserRepo)

	l.Debug("Team service initialized")

	userSvc := service.NewUserService(repo.TeamRepo, repo.UserRepo)

	l.Debug("User service initialized")

	return &Service{
		TeamSvc: teamSvc,
		UserSvc: userSvc,
	}
}

func initHandler(l *zap.Logger, svc *Service) *Handler {
	teamHdl := handler.NewTeamHandler(l, svc.TeamSvc)
	l.Debug("Team handler initialized")

	userHdl := handler.NewUserHandler(l, svc.UserSvc)

	l.Debug("User handler initialized")

	return &Handler{
		TeamHdl: teamHdl,
		UserHdl: userHdl,
	}
}

func initHTTPServer(l *zap.Logger, cfg *config.Config, hdl *Handler) server.HTTPServer {
	router := route.SetupRouter(l, cfg, hdl.TeamHdl, hdl.UserHdl)

	httpServer := server.NewHTTPServer(
		server.WithAddr(cfg.HTTPServer.Host, cfg.HTTPServer.Port),
		server.WithTimeout(cfg.Timeout.Read, cfg.Timeout.Write, cfg.Timeout.Idle),
		server.WithHandler(router),
	)

	return httpServer
}
