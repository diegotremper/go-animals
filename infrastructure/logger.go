package infrastructure

import (
	"os"

	"go.uber.org/zap"
)

func InitLogger() *zap.Logger {
	var (
		logger *zap.Logger
		err    error
	)

	if os.Getenv("APP_ENV") == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic("cannot initialize zap logger: " + err.Error())
	}

	return logger
}
