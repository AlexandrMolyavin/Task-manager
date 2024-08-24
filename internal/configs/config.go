package configs

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct {
	Scheme   string
	Host     string
	Username string
	Password string
	Path     string
}

func InitConfig(logger *zerolog.Logger) *Config {
	viper.AddConfigPath("./configs/")
	viper.SetConfigName("config") // инициализация конфига
	viper.SetConfigType("yaml")

	errcnfg := viper.ReadInConfig()
	if errcnfg != nil {
		logger.Error().Err(errcnfg).Msg("Error reading config file")
		return nil
	}
	var cnfg Config
	cnfg.Scheme = viper.GetString("server.scheme")
	cnfg.Host = viper.GetString("server.localhost")
	cnfg.Username = viper.GetString("server.username")
	cnfg.Password = viper.GetString("server.password")
	cnfg.Path = viper.GetString("server.path")
	return &cnfg

}
