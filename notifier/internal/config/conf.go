package config

import "os"

const (
	// configDirPath путь до папки хранящей конфигурации.
	configDirPath = "config"
	// configType тип разметки файлов конфигурации "json", "toml", "yaml", "yml", "env".
	configType = "env"
	// Dialect типа драйвера/библиотеки при миграции.
	Dialect = "pgx"
)

// AppName имя приложения.
var AppName string

func init() {
	hostname := os.Getenv("HOSTNAME")
	AppName = "notifier"

	if hostname != "" {
		AppName = AppName + "-" + hostname // Значение по умолчанию, если TASK_ID не задано
	}
}
