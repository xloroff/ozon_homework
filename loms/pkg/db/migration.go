package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pressly/goose/v3"

	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/config"
	"gitlab.ozon.dev/xloroff/ozon-hw-go/loms/internal/pkg/logger"
)

/* Запускать миграции из приложения не лучший подход, который несет ряд проблем.
В том числе проблем с разделением прав доступа к БД на миграцию и на обычную работу с базой.
Изначально миграции осуществлялись отдельно, но желание "пощупать" нативную реализацию goose для GO - было важнее. */

// MigrationPool осуществляет миграцию схемы БД.
func MigrationPool(ctx context.Context, l logger.Logger, fldr, cnStr string) error {
	db, err := goose.OpenDBWithDriver(config.Dialect, cnStr)
	if err != nil {
		return fmt.Errorf("Ошибка связи с базой для осуществления миграций - %w", err)
	}

	defer func() {
		err = db.Close()
	}()

	// Проверка версии базы данных.
	currentVersion, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("Ошибка получения версии базы - %w", err)
	}

	err = runMigrations(db, fldr, currentVersion)

	switch {
	case errors.Is(err, goose.ErrNoMigrationFiles):
		l.Warnf(ctx, "db.migrationPool: Миграция не нужна, версия базы %d, задача завершена.", currentVersion)

		err = nil
	case errors.Is(err, nil):
		l.Warnf(ctx, "db.migrationPool: Миграция с версии %d завершена.", currentVersion)
	default:
		l.Errorf(ctx, "db.migrationPool: Ошибка миграции - %v", err)
	}

	return err
}

func runMigrations(db *sql.DB, dir string, currentVersion int64) error {
	// Список всех миграций в папке.
	migrations, err := goose.CollectMigrations(dir, currentVersion+1, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("Ошибка получения списка миграций - %w", err)
	}

	// Применение миграций.
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			if err = goose.UpTo(db, dir, migration.Version); err != nil {
				return fmt.Errorf("Ошибка применения миграции %v - %w", migration.Version, err)
			}
		}
	}

	return nil
}
