package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gomicro-repo/api/config"
	"gomicro-repo/api/logger"
)

func StartHttpServer() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	registerHandlers(engine)

	if err := engine.Run(config.GetConf().Server.Address + config.GetConf().Server.ApiPort); err != nil {
		logger.Error(fmt.Sprintf("StartHttpServer service stop. err: %v", err))
	}
}

func registerHandlers(engine *gin.Engine) {
	engine.POST("/get/userInfo", GetUserInfo)
	engine.POST("/get/all/userInfo", GetAllUserInfo)
	engine.POST("/ping", PingChat)
}
