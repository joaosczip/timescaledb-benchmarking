package configs

import "github.com/spf13/viper"

type Config struct {
	Database struct {
		Host         string `mapstructure:"DB_HOST"`
		Port         int    `mapstructure:"DB_PORT"`
		User         string `mapstructure:"DB_USER"`
		Password     string `mapstructure:"DB_PASSWORD"`
		Name         string `mapstructure:"DB_NAME"`
		SSLModel     string `mapstructure:"DB_SSLMODE"`
		MaxOpenConns int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	}
}

func Load(path string) (*Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
