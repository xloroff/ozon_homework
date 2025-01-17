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
	*BDConSettings
	*JaegerSettings
	*KafkaSettings
}

// AppSettings отвечает за настройки приложения.
type AppSettings struct {
	Port            int     `mapstructure:"API_PORT" validate:"required|int|min:80" message:"Указание порта \"API_PORT\" на котором будет запущен сервис - обязательно"`
	WebPort         int     `mapstructure:"REST_PORT" validate:"required|int|min:80" message:"Указание порта \"REST_PORT\" на котором будет запущен HTTP сервис - обязательно"`
	LogLevel        string  `mapstructure:"LOG_LEVEL" validate:"required|min_len:1" message:"Указание уровня логирования \"LOG_LEVEL\" - обязательно"`
	LogType         int     `mapstructure:"LOG_TYPE" validate:"required|int|min:1|max:3" message:"Необходимо указать тип логирования \"LOG_TYPE\": 0 - логи отключены, 1 - вывод в std.Out, 2 - вывод в файл логов"`
	GracefulTimeout float64 `mapstructure:"GRACEFUL_TIMEOUT"  validate:"float|min:0|max:60" message:"Ожидание закытия соединений с сервисом \"GRACEFUL_TIMEOUT\" не стоит устанавливать более 60 секунд"`
}

// BDConSettings отвечает за настройки подключения к БД.
type BDConSettings struct {
	BDMaster1ConString string `validate:"required|min_len:10" message:"Указание подключения к БД \"DB_NODE_1_CON\" обязательно для передачи через параметры окружения"`
	BDSync1ConString   string `validate:"required|min_len:10" message:"Указание подключения к БД \"DB_SYNC_1_CON\" обязательно для передачи через параметры окружения"`
	MigrationFolder    string `mapstructure:"MIGRATION_FOLDER" validate:"required|min_len:1" message:"Указание папки для миграций  \"MIGRATION_FOLDER\" обязательно для передачи через параметры окружения"`
}

// JaegerSettings отвечает за настройки отправки трейсов.
type JaegerSettings struct {
	JaegerHost string `mapstructure:"JAEGER_HOST" validate:"required|min_len:2" message:"Указание хоста для отправки трейсов \"JAEGER_HOST\" обязательно для передачи через параметры окружения"`
	JaegerPort string `mapstructure:"JAEGER_PORT" validate:"required|min_len:2" message:"Указание порта для отправки трейсов \"JAEGER_PORT\" обязательно для передачи через параметры окружения"`
}

// KafkaSettings отвечает за настройки подключения к Kafka.
type KafkaSettings struct {
	KafkaAddress      string `validate:"required|min_len:4"`
	OrderTopic        string `validate:"required|min_len:1"`
	KafkaSenderPeriod int    `mapstructure:"KAFKA_WAIT_TIMEOUT"  validate:"int|min:1|max:60" message:"Период ожидания доставки сообщения \"KAFKA_WAIT_TIMEOUT\" не стоит устанавливать более 60 секунд"`
	KafkaRetryCount   int    `mapstructure:"KAFKA_RETRY_COUNT" validate:"required|int|min:1|max:10" message:"Количество попыток отправки сообщения\"KAFKA_RETRY_COUNT\" обязательно для передачи"`
	KafkaConnCount    int    `mapstructure:"KAFKA_MAX_CONN" validate:"required|int|min:1|max:10" message:"Количество подключений к Kafka \"KAFKA_MAX_CONN\" обязательно для передачи"`
	KafkaBackoffTime  int    `mapstructure:"KAFKA_BACKOFF_TIME"  validate:"int|min:1|max:600" message:"Период ожидания между попытками отправки (указывается в милисекундах), \"KAFKA_BACKOFF_TIME\" не стоит устанавливать более 6600 миллисекун"`
	BlockTime         int    `mapstructure:"BLOCK_TIME" validate:"required|int|min:1|max:360" message:"Время блокировки сообщений при отправке в брокер \"BLOCK_TIME\" обязательно для передачи. Не более 360 секунд"`
}

// LoadAPIConfig грузит настройки для API.
func LoadAPIConfig() (*ApplicationParameters, error) {
	appSettings, err := loadAppConfig()
	if err != nil {
		return nil, err
	}

	bdSettings, err := loadBDConConfig()
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
	config.BDConSettings = bdSettings
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

// loadBDConConfig грузит настройки подключения к БД.
func loadBDConConfig() (*BDConSettings, error) {
	config := &BDConSettings{}

	config.BDMaster1ConString = os.Getenv("DB_NODE_1_CON")
	config.BDSync1ConString = os.Getenv("DB_SYNC_1_CON")
	config.MigrationFolder = os.Getenv("MIGRATION_FOLDER")

	// Валидируем конфигурацию.
	appBDValidate := validate.Struct(config)
	if !appBDValidate.Validate() {
		err := appBDValidate.Errors
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
	config.OrderTopic = t

	// Валидируем конфигурацию.
	kafkaSettingValidate := validate.Struct(config)
	if !kafkaSettingValidate.Validate() {
		err = kafkaSettingValidate.Errors
		return nil, fmt.Errorf("loadKafkaConfig: ошибка валидации настроек приложения - %w", err)
	}

	return config, nil
}
