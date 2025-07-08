package main

import (
	"github.com/diegotremper/go-animals/config"
	"github.com/diegotremper/go-animals/db"
	"github.com/diegotremper/go-animals/internal/animal"
	"github.com/diegotremper/go-animals/router"
)

func main() {
	config.LoadEnv()
	repo := animal.NewPostgresAnimalRepository(db.InitDB())
	handler := animal.NewAnimalHandler(repo)

	r := router.SetupRouter(handler)
	r.Run()
}
