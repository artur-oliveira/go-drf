package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppName string `mapstructure:"APP_NAME"`

	DBVendor             string `mapstructure:"DB_VENDOR"`
	DBHost               string `mapstructure:"DB_HOST"`
	DBPort               string `mapstructure:"DB_PORT"`
	DBUser               string `mapstructure:"DB_USER"`
	DBPassword           string `mapstructure:"DB_PASSWORD"`
	DBName               string `mapstructure:"DB_NAME"`
	DBLogLevel           string `mapstructure:"DB_LOG_LEVEL"`
	DBMigrate            bool   `mapstructure:"DB_MIGRATE"`
	DBMaxIdle            int    `mapstructure:"DB_MAX_IDLE"`
	DBMaxOpened          int    `mapstructure:"DB_MAX_OPENED"`
	DBMaxLifeTimeSeconds uint   `mapstructure:"DB_MAX_LIFE_TIME_SECONDS"`

	Env string `mapstructure:"ENV"`

	ServerPort         string `mapstructure:"SERVER_PORT"`
	ServerIdleTimeout  int    `mapstructure:"SERVER_IDLE_TIMEOUT"`
	ServerReadTimeout  int    `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout int    `mapstructure:"SERVER_WRITE_TIMEOUT"`

	JWTSecret               string `mapstructure:"JWT_SECRET"`
	JWTExpiresInMinutes     int    `mapstructure:"JWT_EXPIRES_IN_MINUTES"`
	JWTRefreshExpiresInDays int    `mapstructure:"JWT_REFRESH_EXPIRES_IN_DAYS"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.SetDefault("AppName", "GRF")
	viper.SetDefault("DB_VENDOR", "sqlite")
	viper.SetDefault("DB_NAME", "grf")
	viper.SetDefault("DB_LOG_LEVEL", "info")
	viper.SetDefault("DB_MIGRATE", true)
	viper.SetDefault("DB_MAX_IDLE", 10)
	viper.SetDefault("DB_MAX_OPENED", 25)
	viper.SetDefault("DB_MAX_LIFE_TIME_SECONDS", 60)

	viper.SetDefault("ENV", "development") // Padrão seguro

	viper.SetDefault("SERVER_PORT", "1111")
	viper.SetDefault("SERVER_IDLE_TIMEOUT", "3")
	viper.SetDefault("SERVER_READ_TIMEOUT", "3")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "3")

	viper.SetDefault("JWT_SECRET", "my_super_secret_key_insecure_do_not_use_it")
	viper.SetDefault("JWT_EXPIRES_IN_MINUTES", 60*24)
	viper.SetDefault("JWT_REFRESH_EXPIRES_IN_DAYS", 30)

	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
