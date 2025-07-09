package infrastructure

import (
	"time"

	"github.com/diegotremper/go-animals/internal/animal"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(logger *zap.Logger, db *sqlx.DB) *gin.Engine {
	r := gin.Default()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	root := r.Group("/")

	// add module routes
	animal.AddRoutes(logger, db, root)

	return r
}
