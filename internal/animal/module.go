package animal

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Module defines the shared infrastructure dependencies needed by the animal module
type Module interface {
	RootLogger() *slog.Logger
	NewTransactionLogger(ctx *gin.Context) *slog.Logger
	Db() *sqlx.DB
	RouterGroup() *gin.RouterGroup
}
