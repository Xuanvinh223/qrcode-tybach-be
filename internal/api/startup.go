package api

import (
	"fmt"

	router "tyxuan-web-printlabel-api/internal/api/routers"
	"tyxuan-web-printlabel-api/internal/pkg/config"
	"tyxuan-web-printlabel-api/internal/pkg/database"
	"tyxuan-web-printlabel-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Run(configPath string) {
	if configPath == "" {
		configPath = "data/config.yml"
	}
	if err := config.Setup(configPath); err != nil {
		logger.Fatalf("failed to setup config, %s", err)
	}
	if err := database.Setup(); err != nil {
		logger.Fatalf("failed to setup database, %s", err)
	}

	gin.SetMode(config.GetConfig().Server.Mode)

	config := config.GetConfig()

	web := router.Setup()

	fmt.Println("Web API Running on port " + config.Server.Port)
	fmt.Println("================================>")
	logger.Fatalf("%v", web.Run(":"+config.Server.Port))
}
