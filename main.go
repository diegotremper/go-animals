package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// connect db
func init() {
	godotenv.Load()
	dbConn, err := sqlx.Connect("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		panic("database cannot connect")
	}
	db = dbConn
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/animals", createAnimal)
	r.GET("/animals", listAnimals)
	r.GET("/animals/:id", showAnimal)
	r.DELETE("/animals/:id", deleteAnimal)
	r.PATCH("/animals/:id", updateAnimal)

	r.Run()
}

type AnimalCreateRequest struct {
	Name        string
	Age         int
	Description string
}

type AnimalUpdateRequest struct {
	Name        string
	Age         int
	Description string
}

type Animal struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Age         int    `db:"age" json:"age"`
	Description string `db:"description" json:"description"`
}

func createAnimal(ctx *gin.Context) {
	var (
		req          AnimalCreateRequest
		sqlStatement = `
			INSERT INTO animals (name, age, description)
			VALUES ($1, $2, $3)
		`
	)

	ctx.ShouldBindJSON(&req)

	_, err := db.Exec(sqlStatement, req.Name, req.Age, req.Description)
	if err != nil {
		fmt.Println(fmt.Errorf("createAnimal : %w", err))
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "cannot create animal"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success create animal"})
}

func listAnimals(ctx *gin.Context) {
	var (
		animals      []Animal
		sqlStatement = `
			SELECT id, name, age, description
			FROM animals
		`
	)

	rows, err := db.Queryx(sqlStatement)
	if err != nil {
		fmt.Println(fmt.Errorf("listAnimals : %w", err))
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Cannot get list animals"})
		return
	}

	for rows.Next() {
		var animal Animal
		rows.StructScan(&animal)
		animals = append(animals, animal)
	}

	ctx.JSON(http.StatusOK, animals)
}

func showAnimal(ctx *gin.Context) {
	var (
		id, _        = ctx.Params.Get("id")
		animal       Animal
		sqlStatement = `
			SELECT *
			FROM animals
			WHERE id = $1
		`
	)

	err := db.QueryRowx(sqlStatement, id).StructScan(&animal)
	if err != nil {
		fmt.Println(fmt.Errorf("showAnimal : %w", err))
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Animal not found"})
		return
	}

	ctx.JSON(http.StatusOK, animal)
}

func deleteAnimal(ctx *gin.Context) {
	var (
		id, _        = ctx.Params.Get("id")
		sqlStatement = `
			DELETE
			FROM animals
			WHERE id = $1
		`
	)

	result, err := db.Exec(sqlStatement, id)
	if err != nil {
		fmt.Println(fmt.Errorf("deleteAnimal : %w", err))
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Animal not found"})
		return
	}

	count, _ := result.RowsAffected()
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("delete %d animal record.", count)})
}

func updateAnimal(ctx *gin.Context) {
	var (
		req          AnimalUpdateRequest
		id, _        = ctx.Params.Get("id")
		sqlStatement = `
			UPDATE animals
			SET name = $1, age = $2, description = $3
			WHERE id = $4
		`
	)

	ctx.ShouldBindJSON(&req)

	result, err := db.Exec(sqlStatement, req.Name, req.Age, req.Description, id)
	if err != nil {
		fmt.Println(fmt.Errorf("updateAnimal : %w", err))
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Cannot update animal."})
		return
	}

	count, _ := result.RowsAffected()
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("update %d animal record.", count)})
}
