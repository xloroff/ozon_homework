package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gookit/validate"
	"github.com/spf13/viper"
)

// ApplicationParameters конфиг для API/приложения и коммуникации с внешними сервисами.
type ApplicationParameters struct {
	*AppSettings
	*JaegerSettings
	*KafkaSettings
}

// AppSettings отвечает за настройки приложения.
type AppSettings struct {
	LogLevel        string  `mapstructure:"LOG_LEVEL" validate:"required|min_len:1" message:"Указание уровня логирования \"LOG_LEVEL\" - обязательно"`
	LogType         int     `mapstructure:"LOG_TYPE" validate:"required|int|min:1|max:3" message:"Необходимо указать тип логирования \"LOG_TYPE\": 0 - логи отключены, 1 - вывод в std.Out, 2 - вывод в файл логов"`
	GracefulTimeout float64 `mapstructure:"GRACEFUL_TIMEOUT"  validate:"float|min:0|max:60" message:"Ожидание закытия соединений с сервисом \"GRACEFUL_TIMEOUT\" не стоит устанавливать более 60 секунд"`
}

// JaegerSettings отвечает за настройки отправки трейсов.
type JaegerSettings struct {
	JaegerHost string `mapstructure:"JAEGER_HOST" validate:"required|min_len:2" message:"Указание хоста для отправки трейсов \"JAEGER_HOST\" обязательно для передачи через параметры окружения"`
	JaegerPort string `mapstructure:"JAEGER_PORT" validate:"required|min_len:2" message:"Указание порта для отправки трейсов \"JAEGER_PORT\" обязательно для передачи через параметры окружения"`
}

// KafkaSettings отвечает за настройки подключения к Kafka.
type KafkaSettings struct {
	KafkaAddress           string `validate:"required|min_len:4"`
	KafkaTopic             string `validate:"required|min_len:1"`
	KafkaConsumerGroupName string `validate:"required|min_len:1"`
	KafkaHeartbeatTime     int    `mapstructure:"KAFKA_HEARTBEAT_TIME"  validate:"int|min:1|max:600" message:"Подача сигналов признаков жизни констюмера \"KAFKA_HEARTBEAT_TIME\", обязательна для передачи - не стоит устанавливать более 600 секунд"`
	KafkaSessionTimeout    int    `mapstructure:"KAFKA_SESSION_TIMEOUT" validate:"required|int|min:1|max:600" message:"Таймаут сессии с Kafka \"KAFKA_SESSION_TIMEOUT\" обязателен для передачи - не стоит устанавливать более 600 секунд"`
	KafkaRebalanceTime     int    `mapstructure:"KAFKA_REBALANCE_TIME" validate:"required|int|min:1|max:600" message:"Время таймаута для ребаланса косьюмеров Kafka \"KAFKA_REBALANCE_TIME\" обязательно для передачи - не стоит устанавливать более 600 секунд"`
	KafkaAutoCommitTime    int    `mapstructure:"KAFKA_AUTOCOMMIT_TIME"  validate:"int|min:1|max:600" message:"Интервал автоматических коммитов смещения в Kafka (указывается в секундах), \"KAFKA_AUTOCOMMIT_TIME\" обязателен для передачи (даже если автокоммит отключен) - не стоит устанавливать более 600 секунд"`
	KafkaAutocommitEnabled bool   `mapstructure:"KAFKA_AUTO_COMMIT"  validate:"isBool" message:"Параметр определения автоматической фиксации смещения сообщений \"KAFKA_AUTO_COMMIT\" обязателен для передачи"`
}

// LoadAPIConfig грузит настройки для API.
func LoadAPIConfig() (*ApplicationParameters, error) {
	appSettings, err := loadAppConfig()
	if err != nil {
		return nil, err
	}

	jaegerSettings, err := loadJaegerConfig()
	if err != nil {
		return nil, err
	}

	kafkaSettings, err := loadKafkaConfig()
	if err != nil {
		return nil, err
	}

	config := &ApplicationParameters{}
	config.AppSettings = appSettings
	config.JaegerSettings = jaegerSettings
	config.KafkaSettings = kafkaSettings

	return config, nil
}

func loadAppConfig() (*AppSettings, error) {
	// Файлики конфигов.
	viper.AddConfigPath(configDirPath)
	viper.SetConfigType(configType)
	viper.SetConfigName(os.Getenv("APP_CONFIG_NAME"))

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

// loadJaegerConfig грузит настройки подключения к Jaeger.
func loadJaegerConfig() (*JaegerSettings, error) {
	config := &JaegerSettings{}

	config.JaegerHost = os.Getenv("JAEGER_HOST")
	config.JaegerPort = os.Getenv("JAEGER_PORT")

	// Валидируем конфигурацию.
	appJaegerValidate := validate.Struct(config)
	if !appJaegerValidate.Validate() {
		err := appJaegerValidate.Errors
		return nil, fmt.Errorf("LoadJaegerConfig: ошибка валидации настроек связи c Jaeger - %w", err)
	}

	return config, nil
}

// loadKafkaConfig грузит настройки подключения к Kafka.
func loadKafkaConfig() (*KafkaSettings, error) {
	viper.AddConfigPath(configDirPath)
	viper.SetConfigType(configType)
	viper.AutomaticEnv()

	config := &KafkaSettings{}

	h := os.Getenv("KAFKA_HOST")
	p := os.Getenv("KAFKA_PORT")
	t := os.Getenv("TOPIC_NAME")
	n := os.Getenv("CONSUMER_GROUP_NAME")

	if h == "" || p == "" {
		return nil, fmt.Errorf("Ошибка получения данных подключения к Kafka - хост \"%v\", порт \"%v\"", h, p)
	}

	urlString := fmt.Sprintf("%v:%v", h, p)

	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, fmt.Errorf("Ошибка форматирования данных подключения к Kafka - %w", err)
	}

	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("loadKafkaConfig: ошибка чтения конфигурации из файла - %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("loadKafkaConfig: ошибка преобразования переменных из конфигурации - %w", err)
	}

	config.KafkaAddress = urlString
	config.KafkaTopic = t
	config.KafkaConsumerGroupName = n

	// Валидируем конфигурацию.
	kafkaSettingValidate := validate.Struct(config)
	if !kafkaSettingValidate.Validate() {
		err = kafkaSettingValidate.Errors
		return nil, fmt.Errorf("loadKafkaConfig: ошибка валидации настроек приложения - %w", err)
	}

	return config, nil
}
