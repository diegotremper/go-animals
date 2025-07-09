package animal

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func AddRoutes(log *zap.Logger, db *sqlx.DB, rg *gin.RouterGroup) {
	animals := rg.Group("/animals")

	repo := NewPostgresAnimalRepository(log, db)
	handler := NewAnimalHandler(repo)

	animals.POST("", handler.CreateAnimalHandler)
	animals.GET("/:id", handler.GetAnimalHandler)
	animals.GET("", handler.ListAnimalsHandler)
	animals.PUT("/:id", handler.UpdateAnimalHandler)
	animals.DELETE("/:id", handler.DeleteAnimalHandler)
}
