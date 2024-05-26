package initilize

import (
	"fmt"

	"github.com/gookit/validate"
	"github.com/spf13/viper"
)

// ConfigAPI конфиг для API/приложения и коммуникации с внешними сервисами.
type ConfigAPI struct {
	Port                int     `mapstructure:"API_PORT" validate:"required|int|min:80" message:"Указание порта \"API_PORT\" на котором будет запущен сервис - обязательно"`
	LogLevel            string  `mapstructure:"LOG_LEVEL" validate:"required|min_len:1" message:"Указание уровня логирования \"LOG_LEVEL\" - обязательно"`
	LogFolder           string  `mapstructure:"LOG_FOLDER" validate:"required|unixPath" message:"Необходимо указать полный путь до папки хранения логов \"LOG_FOLDER\""`
	ProductServiceHost  string  `mapstructure:"PRODUCT_SERVICE_HOST" validate:"required|fullUrl" message:"Адрес сервиса продуктов \"PRODUCT_SERVICE_HOST\" необходимо указать полностью включая http или https"`
	ProductServiceToken string  `mapstructure:"PRODUCT_SERVICE_TOKEN" validate:"required|min_len:1" message:"Необходимо указать токен для авторизации на сервисе продуктов \"PRODUCT_SERVICE_TOKEN\""`
	ProductServiceRetr  int     `mapstructure:"PRODUCT_SERVICE_RETRIES" validate:"required|int|min:1|max:20" message:"Число попыток установки связи с сервисом продуктов \"PRODUCT_SERVICE_RETRIES\" должно быть от 1 до 20"`
	GracefulTimeout     float64 `mapstructure:"GRACEFUL_TIMEOUT"  validate:"float|min:0|max:60" message:"Ожидание закытия соединений с сервисом \"GRACEFUL_TIMEOUT\" не стоит устанавливать более 60 секунд"`
	ClientVer           int     `mapstructure:"CLIENT_PRODUCT_SERVICE" validate:"required|int|min:1|max:2" message:"Версию используемого клиента \"CLIENT_PRODUCT_SERVICE\" можно назначить только из набора: 1 - стандартная библиотека 2 - go-resty"`
}

// LoadApiConfig грузит настройки для API.
func LoadApiConfig(configDirPath, configType, appConfigName string) (ConfigAPI, error) {
	// Файлики конфигов.
	viper.AddConfigPath(configDirPath)
	viper.SetConfigType(configType)
	viper.SetConfigName(appConfigName)

	viper.AutomaticEnv()

	config := ConfigAPI{}

	// Читаем данные и раскидываем.
	err := viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("Ошибка чтения конфигурационного файла - %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("Ошибка преобразования конфигурационного файла - %w", err)
	}

	// Валидируем конфигурацию.
	v := validate.Struct(config)
	if !v.Validate() {
		err = v.Errors
		return config, fmt.Errorf("Ошибка валидации параметров конфигурационного файла - %w", err)
	}

	return config, nil
}