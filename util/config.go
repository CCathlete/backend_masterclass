package u

import "github.com/spf13/viper"

// States all configurations of the application.
// The values are read by viper from a .env file.
type Config struct {
	Host          string `mapstructure:"DB_HOST"`
	Port          string `mapstructure:"DB_PORT"`
	User          string `mapstructure:"DB_USER"`
	Password      string `mapstructure:"DB_PASS"`
	DBName        string `mapstructure:"DB_NAME"`
	SSLMode       string `mapstructure:"DB_SSL_MODE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// Reads configuration from files or env variables.
func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Reading the environment variables into viper's buffer.
	viper.AutomaticEnv()

	// Reading the .env file into viper's buffer.
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
