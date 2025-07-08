package router

import (
	"github.com/diegotremper/go-animals/internal/animal"

	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *animal.AnimalHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/animals", handler.CreateAnimalHandler)
	r.GET("/animals/:id", handler.GetAnimalHandler)
	r.GET("/animals", handler.ListAnimalsHandler)
	r.PUT("/animals/:id", handler.UpdateAnimalHandler)
	r.DELETE("/animals/:id", handler.DeleteAnimalHandler)

	return r
}
