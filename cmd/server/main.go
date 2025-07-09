package main

import (
	"github.com/diegotremper/go-animals/infrastructure"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	logger := infrastructure.InitLogger()
	db := infrastructure.InitDB()
	r := infrastructure.SetupRouter(logger, db)

	r.Run()
}
