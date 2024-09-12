package util

import "github.com/spf13/viper"

type Config struct {
	DBHost         string `mapstructure:"MYSQL_HOST"`
	DBUserName     string `mapstructure:"MYSQL_USER"`
	DBUserPassword string `mapstructure:"MYSQL_PASSWORD"`
	DBName         string `mapstructure:"MYSQL_DB"`
	DBPort         string `mapstructure:"MYSQL_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	S3BucketName string `mapstructure:"S3BucketName"`
	S3Region     string `mapstructure:"S3Region"`
	S3APIKey     string `mapstructure:"S3APIKey"`
	S3SecretKey  string `mapstructure:"S3SecretKey"`
	S3Domain     string `mapstructure:"S3Domain"`

	SecretKey string `mapstructure:"SYSTEM_SECRET"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	// viper.SetConfigName("app")
	// viper.SetConfigType(".env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
