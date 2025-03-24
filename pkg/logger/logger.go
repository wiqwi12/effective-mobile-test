package logger

import (
	"fmt"
	"github.com/wiqwi12/effective-mobile-test/pkg/cfg"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type Logger struct {
	Debug *slog.Logger
	Info  *slog.Logger
}

func NewLogger(cfg cfg.Config) (*Logger, error) {

	if cfg.DebugFilePath == "" {
		cfg.DebugFilePath = "logs/debug.log"
	}

	if cfg.InfoFilePath == "" {
		cfg.InfoFilePath = "logs/info.log"
	}

	var debugLogger, infoLogger *slog.Logger

	//consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	//	Level: slog.LevelDebug,
	//})

	// Create debug log file and handler
	debugDir := filepath.Dir(cfg.DebugFilePath)
	if err := os.MkdirAll(debugDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create debug log directory %s: %w", debugDir, err)
	}

	debugFile, err := os.OpenFile(cfg.DebugFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open debug log file %s: %w", cfg.DebugFilePath, err)
	}

	debugHandler := slog.NewJSONHandler(debugFile, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	debugLogger = slog.New(debugHandler)

	infoDir := filepath.Dir(cfg.InfoFilePath)
	if err := os.MkdirAll(infoDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create info log directory %s: %w", infoDir, err)
	}

	infoFile, err := os.OpenFile(cfg.InfoFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open info log file %s: %w", cfg.InfoFilePath, err)
	}

	infoHandler := slog.NewJSONHandler(infoFile, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	infoLogger = slog.New(infoHandler)

	if cfg.ConsoleOutput == "true" {
		debugMultiWriter := io.MultiWriter(debugFile, os.Stdout)
		debugMultiHandler := slog.NewJSONHandler(debugMultiWriter, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
		debugLogger = slog.New(debugMultiHandler)

		infoMultiWriter := io.MultiWriter(infoFile, os.Stdout)
		infoMultiHandler := slog.NewJSONHandler(infoMultiWriter, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
		infoLogger = slog.New(infoMultiHandler)
	}

	return &Logger{
		Debug: debugLogger,
		Info:  infoLogger,
	}, nil
}
