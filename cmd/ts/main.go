package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	cons2 "ts/internal/broker/kafka/cons"
	prod2 "ts/internal/broker/kafka/prod"
	"ts/internal/configs"
	"ts/internal/controllers"
	"ts/internal/repository/mdb"
	"ts/internal/repository/pdb"
	"ts/internal/transport"
	"ts/pkg/broker/kafka"
	"ts/pkg/logger"
	"ts/pkg/middleware"
	"ts/pkg/repository/mongoDB"
	"ts/pkg/repository/postgresql"
)

const port string = ":8081"

func main() {

	router := gin.Default()

	zlog := logger.InitLogger()      // Logger Init
	cnfg := configs.InitConfig(zlog) // Get info from config.yaml

	var handler *transport.Handler     // Task Handler
	var uHandler *controllers.UHandler // User Handler

	switch cnfg.Scheme {
	case "postgres":
		db := postgresql.GetConnection(zlog, cnfg)
		repo := pdb.NewPgx(db)
		handler = transport.NewHandler(&repo)
		uHandler = controllers.NewHandler(&repo)
		handler.Logger, uHandler.Logger = zlog, zlog
	case "mongodb":
		db := mongoDB.GetConnection(zlog)
		repo := mdb.NewMongoDB(db)
		handler = transport.NewHandler(&repo)
	}

	addr := []string{"localhost:9092"}
	prod, errP := kafka.ConnectProducer(addr)
	if errP != nil {
		zlog.Fatal().Err(errP).Msg("Error starting Producer")
		return
	}
	sProd := prod2.NewProducer(prod)
	cons, errC := kafka.ConnectConsumer(addr)
	if errC != nil {
		zlog.Fatal().Err(errC).Msg("Error starting Consumer")
		return
	}
	sCons := cons2.NewConsumer(cons)

	handler.Prod = sProd
	handler.Cons = sCons

	router.POST("/signup", controllers.SignUp(uHandler)) // Registration
	router.POST("/login", controllers.Login(uHandler))   // Login

	router.Any("/", middleware.RequireAuth(uHandler), handler.HandleRequest)

	zlog.Info().Msg("Starting server...")
	str := fmt.Sprintf("Listening on port%s", port)
	zlog.Info().Msg(str)
	err := router.Run(port)

	if err != nil {
		zlog.Fatal().Err(err).Msg("Error starting server")
		return
	}
}
