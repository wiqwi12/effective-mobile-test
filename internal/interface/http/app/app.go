// @title Music Library API
// @version 1.0
// @description API для управления библиотекой песен и исполнителей

// @host localhost:8080
// @BasePath /api
// @schemes http

package app

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
	"github.com/wiqwi12/effective-mobile-test/internal/infrastructure/externalServices"
	"github.com/wiqwi12/effective-mobile-test/internal/infrastructure/postgres/migration"
	"github.com/wiqwi12/effective-mobile-test/internal/infrastructure/postgres/repository"
	"github.com/wiqwi12/effective-mobile-test/internal/interface/http/handlers"
	"github.com/wiqwi12/effective-mobile-test/internal/interface/http/middleware"
	"github.com/wiqwi12/effective-mobile-test/internal/service"
	"github.com/wiqwi12/effective-mobile-test/pkg"
	"github.com/wiqwi12/effective-mobile-test/pkg/cfg"
	"github.com/wiqwi12/effective-mobile-test/pkg/logger"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	//nolint
	_ "github.com/wiqwi12/effective-mobile-test/docs"
)

func Run() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	workdir, err := os.Getwd()
	if err != nil {
		fmt.Errorf("failed to get working directory: %w", err)
	}

	migrPath := filepath.Join(workdir, "internal/infrastructure/postgres/migration")

	psqlCfg := cfg.PSQLconfig{
		Host:          os.Getenv("POSTGRES_HOST"),
		Port:          os.Getenv("POSTGRES_PORT"),
		Username:      os.Getenv("POSTGRES_USER"),
		Password:      os.Getenv("POSTGRES_PASSWORD"),
		Database:      os.Getenv("POSTGRES_DB"),
		MigrationPath: migrPath,
	}

	httpConfig := cfg.HTTPconfig{
		Host: os.Getenv("HTTP_HOST"),
		Port: os.Getenv("HTTP_PORT"),
	}
	if httpConfig.Host == "" && httpConfig.Port == "" {
		httpConfig.Port = "8080"
		httpConfig.Host = "localhost"
	}

	loggerConfig := cfg.Config{
		DebugFilePath: os.Getenv("DEBUG_FILE_PATH"),
		InfoFilePath:  os.Getenv("INFO_FILE_PATH"),
		ConsoleOutput: os.Getenv("CONSOLE_OUTPUT"),
	}

	logger, err := logger.NewLogger(loggerConfig)
	if err != nil {
		slog.Error(err.Error())
	}

	externalServiceApi := os.Getenv("EXTERNAL_SERVICE_API")

	db, err := pkg.NewDbConn(psqlCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = migration.RunMigrations(db)
	if err != nil {
		slog.Error(err.Error())
	}
	logger.Info.Info("migrations applied succsessfully")

	groupRepo := repository.NewGroupRepository(db, logger)
	songRepo := repository.NewSongRepo(db, logger)
	MetadataRepo := externalServices.NewExternalRepo(externalServiceApi, logger)
	verseRepo := repository.NewVerseRepository(db, logger)
	service := service.NewSongSrvc(songRepo, groupRepo, MetadataRepo, verseRepo, logger)
	validator := validator.New()

	handler := handlers.NewHandler(service, validator)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/song", handler.CreateSongHandler)
	mux.HandleFunc("GET /api/song/{id}", handler.GetSongHandler)
	mux.HandleFunc("PUT /api/song/{id}", handler.UpdateSongHandler)
	mux.HandleFunc("DELETE /api/song/{id}", handler.DeleteSongHandler)
	mux.HandleFunc("GET /api/song", handler.GetSongWithFilter)
	mux.HandleFunc("GET /api/verses/{id}", handler.GetPaginatedVerses)
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	headersMWMux := middleware.CommonHeadersMiddleware(mux)

	server := &http.Server{
		Addr:    httpConfig.Host + ":" + httpConfig.Port,
		Handler: headersMWMux,
	}

	go func() {
		logger.Debug.Info("Starting Server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Debug.Info("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %s", err)
	}

	logger.Debug.Info("Server gracefully stopped")
}
