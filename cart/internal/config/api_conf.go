package config

import (
	"fmt"

	"github.com/gookit/validate"
	"github.com/spf13/viper"
)

// ApplicationParameters конфиг для API/приложения и коммуникации с внешними сервисами.
type ApplicationParameters struct {
	*AppSettings
	*ProductServiceSettings
}

// AppSettings отвечает за настройки приложения.
type AppSettings struct {
	Port            int     `mapstructure:"API_PORT" validate:"required|int|min:80" message:"Указание порта \"API_PORT\" на котором будет запущен сервис - обязательно"`
	LogLevel        string  `mapstructure:"LOG_LEVEL" validate:"required|min_len:1" message:"Указание уровня логирования \"LOG_LEVEL\" - обязательно"`
	LogType         int     `mapstructure:"LOG_TYPE" validate:"required|int|min:1|max:3" message:"Необходимо указать тип логирования \"LOG_TYPE\": 0 - логи отключены, 1 - вывод в std.Out, 2 - вывод в файл логов"`
	GracefulTimeout float64 `mapstructure:"GRACEFUL_TIMEOUT"  validate:"float|min:0|max:60" message:"Ожидание закытия соединений с сервисом \"GRACEFUL_TIMEOUT\" не стоит устанавливать более 60 секунд"`
}

// ProductServiceSettings отвечает за настройки связи с сервисом продуктов.
type ProductServiceSettings struct {
	ProductServiceHost  string `mapstructure:"PRODUCT_SERVICE_HOST" validate:"required|fullUrl" message:"Адрес сервиса продуктов \"PRODUCT_SERVICE_HOST\" необходимо указать полностью включая http или https"`
	ProductServiceToken string `mapstructure:"PRODUCT_SERVICE_TOKEN" validate:"required|min_len:1" message:"Необходимо указать токен для авторизации на сервисе продуктов \"PRODUCT_SERVICE_TOKEN\""`
	ProductServiceRetr  int    `mapstructure:"PRODUCT_SERVICE_RETRIES" validate:"required|int|min:1|max:20" message:"Число попыток установки связи с сервисом продуктов \"PRODUCT_SERVICE_RETRIES\" должно быть от 1 до 20"`
}

// LoadAPIConfig грузит настройки для API.
func LoadAPIConfig() (*ApplicationParameters, error) {
	appSettings, err := loadAppConfig()
	if err != nil {
		return nil, err
	}

	productSettings, err := loadProductSettings()
	if err != nil {
		return nil, err
	}

	config := &ApplicationParameters{}
	config.AppSettings = appSettings
	config.ProductServiceSettings = productSettings

	return config, nil
}

func loadAppConfig() (*AppSettings, error) {
	// Файлики конфигов.
	viper.AddConfigPath(configDirPath)
	viper.SetConfigType(configType)
	viper.SetConfigName(appConfigName)

	viper.AutomaticEnv()

	config := &AppSettings{}

	// Читаем данные и раскидываем.
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("LoadAPIConfig: ошибка чтения конфигурации из файла - %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("LoadAPIConfig: ошибка преобразования переменных из конфигурации - %w", err)
	}

	// Валидируем конфигурацию.
	appSettingValidate := validate.Struct(config)
	if !appSettingValidate.Validate() {
		err = appSettingValidate.Errors
		return nil, fmt.Errorf("LoadAPIConfig: ошибка валидации настроек приложения - %w", err)
	}

	return config, nil
}

func loadProductSettings() (*ProductServiceSettings, error) {
	// Файлики конфигов.
	viper.AddConfigPath(configDirPath)
	viper.SetConfigType(configType)
	viper.SetConfigName(appConfigName)

	viper.AutomaticEnv()

	config := &ProductServiceSettings{}

	// Читаем данные и раскидываем.
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("LoadAPIConfig: ошибка чтения конфигурации из файла - %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("LoadAPIConfig: ошибка преобразования переменных из конфигурации - %w", err)
	}

	// Валидируем конфигурацию.
	appSettingValidate := validate.Struct(config)
	if !appSettingValidate.Validate() {
		err = appSettingValidate.Errors
		return nil, fmt.Errorf("LoadAPIConfig: ошибка валидации параметров связи с сервисом продуктов - %w", err)
	}

	return config, nil
}
