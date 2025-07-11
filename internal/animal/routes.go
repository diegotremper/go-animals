package animal

func AddRoutes(module Module) {
	rg := module.RouterGroup()
	db := module.Db()

	animals := rg.Group("/animals")

	repo := NewPostgresAnimalRepository(db)
	handler := NewAnimalHandler(module, repo)

	animals.POST("", handler.CreateAnimalHandler)
	animals.GET("/:id", handler.GetAnimalHandler)
	animals.GET("", handler.ListAnimalsHandler)
	animals.PUT("/:id", handler.UpdateAnimalHandler)
	animals.DELETE("/:id", handler.DeleteAnimalHandler)
}
