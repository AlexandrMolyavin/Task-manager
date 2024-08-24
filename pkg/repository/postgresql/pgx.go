package postgresql

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
	"ts/internal/configs"
)

//var Pg = NewPgx() // Postgres

type Config struct {
	scheme   string
	host     string
	username string
	password string
	path     string
}

func initDsn() *Config {
	var cnfg Config
	cnfg.scheme = viper.GetString("server.scheme")
	cnfg.host = viper.GetString("server.localhost")
	cnfg.username = viper.GetString("server.username")
	cnfg.password = viper.GetString("server.password")
	cnfg.path = viper.GetString("server.path")
	return &cnfg
}

func GetConnection(logger *zerolog.Logger, cnfg *configs.Config) *gorm.DB {

	logger.Info().Msg("Connecting to database...")
	dsn := url.URL{
		Scheme: cnfg.Scheme,
		Host:   cnfg.Host, // postgresql:5432 для контейнера, localhost для пк, 127.0.0.1:15432
		User:   url.UserPassword(cnfg.Username, cnfg.Password),
		Path:   cnfg.Path,
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")
	dsn.RawQuery = q.Encode()
	dsnStr := dsn.String()

	db, err := gorm.Open(postgres.Open(dsnStr), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("Cannot connect to postgres")
		panic(err)
		return nil
	}
	logger.Info().Msg("Connection Successful")

	return db
}
