package config

const (
	// ConfigDirPath путь до папки хранящей конфигурации.
	ConfigDirPath = "../config"
	// ConfigType тип разметки файлов конфигурации "json", "toml", "yaml", "yml", "env".
	ConfigType = "env"
	// AppConfigName имя файлика конфигурации.
	AppConfigName = "api.conf"
)

const (
	// LogfileName имя файла с логами.
	LogfileName = "/app.log"
)

// HTTPCalmStatus Какая-то нестандартная ошибка которая может появляться при коммуникации с внешним сервисом продуктов.
const HTTPCalmStatus = 420