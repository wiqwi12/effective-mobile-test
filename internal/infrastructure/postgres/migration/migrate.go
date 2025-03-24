package migration

import (
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pressly/goose/v3"
	"os"
	"path/filepath"
	"runtime"
)

func RunMigrations(db *sql.DB) error {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}

	migrationsPath := filepath.Dir(currentFilePath)
	fmt.Println(migrationsPath)

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationsPath)
	}

	goose.SetDialect("postgres")

	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}
