package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/portfolio/backend/internal/config"
	"github.com/portfolio/backend/internal/database"
	"github.com/portfolio/backend/internal/handler"
	"github.com/portfolio/backend/internal/router"
	"github.com/portfolio/backend/internal/service"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	if err := run(log); err != nil {
		log.Error("server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log.Info("config loaded", slog.String("port", cfg.Port))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := database.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Info("database connected")

	authSvc := service.NewAuthService(cfg.JWTSecret, cfg.AdminPasswordHash)
	projectSvc := service.NewProjectService(db)
	contentSvc := service.NewContentService(db)
	careerSvc := service.NewCareerService(db)
	inquirySvc := service.NewInquiryService(db)
	uploadSvc := service.NewUploadService(db, cfg.UploadDir)

	cacheTTL := time.Duration(cfg.UmamiCacheTTLSeconds) * time.Second
	analyticsSvc, err := service.NewAnalyticsService(
		cfg.UmamiAPIURL,
		cfg.UmamiAPIKey,
		cfg.UmamiWebsiteID,
		cfg.AppEnv,
		cacheTTL,
		log,
	)
	if err != nil {
		return err
	}

	handlers := router.Handlers{
		Health:    handler.NewHealthHandler(db),
		Auth:      handler.NewAuthHandler(authSvc),
		Project:   handler.NewProjectHandler(projectSvc),
		Content:   handler.NewContentHandler(contentSvc),
		Career:    handler.NewCareerHandler(careerSvc),
		Inquiry:   handler.NewInquiryHandler(inquirySvc),
		Upload:    handler.NewUploadHandler(uploadSvc),
		Analytics: handler.NewAnalyticsHandler(analyticsSvc, log),
	}

	r := router.New(router.Deps{
		Ctx:         ctx,
		Handlers:    handlers,
		AuthService: authSvc,
		CORSOrigins: cfg.CORSOrigins,
		Logger:      log,
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info("server starting", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", slog.String("error", err.Error()))
			cancel()
		}
	}()

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Info("server stopped")
	return nil
}
