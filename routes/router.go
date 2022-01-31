package routes

import (
	"wordAvailable/controllers"
	_ "wordAvailable/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// InitRouter init router
func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if mode := gin.Mode(); mode == gin.DebugMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	router.GET("/heartBeat", controllers.HeartBeat)

	apiv1 := router.Group("/api/v1")
	apiv1.GET("/words", controllers.GetWords)

	return router
}
