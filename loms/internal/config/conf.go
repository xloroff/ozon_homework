package config

const (
	// configDirPath путь до папки хранящей конфигурации.
	configDirPath = "config"
	// configType тип разметки файлов конфигурации "json", "toml", "yaml", "yml", "env".
	configType = "env"
)

// HTTPCalmStatus Какая-то нестандартная ошибка которая может появляться при коммуникации с внешним сервисом продуктов.
const HTTPCalmStatus = 420
