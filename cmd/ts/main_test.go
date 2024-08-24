package main

import (
	"github.com/gin-gonic/gin"
	"testing"
	"ts/internal/configs"
	"ts/internal/controllers"
	"ts/internal/repository/pdb"
	"ts/pkg/logger"
	"ts/pkg/repository/postgresql"
)

func TestGormDB(t *testing.T) {
	router := gin.Default()
	var zlog = logger.InitLogger()   // Logger Init
	cnfg := configs.InitConfig(zlog) // Get info from config.yaml
	db := postgresql.GetConnection(zlog, cnfg)
	repo := pdb.NewPgx(db)
	var uHandler *controllers.UHandler
	uHandler = controllers.NewHandler(&repo)
	router.POST("/login", controllers.Login(uHandler)) // Login
}
