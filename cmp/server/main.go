package main

import (
	"github.com/diegotremper/go-animals/config"
	"github.com/diegotremper/go-animals/db"
	"github.com/diegotremper/go-animals/router"
)

func main() {
	config.LoadEnv()
	db.InitDB()
	r := router.SetupRouter()
	r.Run()
}
