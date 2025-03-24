package pkg

import (
	"database/sql"
	"fmt"
	"github.com/wiqwi12/effective-mobile-test/pkg/cfg"
	"log"
	//nolint
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDbConn(cfg cfg.PSQLconfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	fmt.Println(dsn)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}

	log.Println("Подключение к PostgreSQL успешно")
	return db, nil
}
