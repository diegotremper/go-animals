package infrastructure

import (
	"log/slog"

	"github.com/diegotremper/go-animals/internal/animal"
	sloggin "github.com/samber/slog-gin"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AnimalModule struct {
	logger *slog.Logger
	db     *sqlx.DB
	rg     *gin.RouterGroup
}

func (m *AnimalModule) RootLogger() *slog.Logger {
	return m.logger
}

func (m *AnimalModule) NewTransactionLogger(ctx *gin.Context) *slog.Logger {
	return NewTransactionLogger(ctx)
}

func (m *AnimalModule) Db() *sqlx.DB {
	return m.db
}

func (m *AnimalModule) RouterGroup() *gin.RouterGroup {
	return m.rg
}

func SetupRouter(logger *slog.Logger, db *sqlx.DB) *gin.Engine {
	r := gin.Default()

	r.Use(sloggin.New(logger))
	r.Use(gin.Recovery())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	root := r.Group("/")

	// add module routes
	animal.AddRoutes(&AnimalModule{
		logger: logger,
		db:     db,
		rg:     root,
	})

	return r
}
