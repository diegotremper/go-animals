package router

import (
	"github.com/diegotremper/go-animals/internal/animal"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/animals", animal.CreateAnimalHandler)
	r.GET("/animals/:id", animal.GetAnimalHandler)
	r.GET("/animals", animal.ListAnimalsHandler)
	r.PUT("/animals/:id", animal.UpdateAnimalHandler)
	r.DELETE("/animals/:id", animal.DeleteAnimalHandler)

	return r
}
