package animal

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/diegotremper/go-animals/db"
)

type AnimalRepository interface {
	CreateAnimal(r AninalCreateRequest) error
	UpdateAnimal(id int, r AnimalUpdateRequest) error
	ListAnimals() ([]Animal, error)
	GetAnimal(id int) (Animal, error)
	DeleteAnimal(id int) error
}

type PostgresAnimalRepository struct{}

func (PostgresAnimalRepository) CreateAnimal(r AninalCreateRequest) error {
	_, err := db.DB.Exec(`INSERT INTO animals (name, age, description) VALUES ($1, $2, $3)`, r.Name, r.Age, r.Description)

	return err
}

func (PostgresAnimalRepository) UpdateAnimal(id int, r AnimalUpdateRequest) error {
	_, err := db.DB.Exec(`UPDATE animals SET name = $1, age = $2, description = $3 WHERE id = $4`, r.Name, r.Age, r.Description, id)

	return err
}

func (PostgresAnimalRepository) ListAnimals() ([]Animal, error) {
	var (
		animals      []Animal = make([]Animal, 0)
		sqlStatement          = `
			SELECT id, name, age, description FROM animals
		`
	)

	rows, err := db.DB.Queryx(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("ListAnimals query error: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var animal Animal
		if err := rows.StructScan(&animal); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		animals = append(animals, animal)
	}

	return animals, nil
}

func (PostgresAnimalRepository) GetAnimal(id int) (Animal, error) {
	var (
		animal       Animal
		sqlStatement = `SELECT * FROM animals WHERE id = $1`
	)

	err := db.DB.QueryRowx(sqlStatement, id).StructScan(&animal)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Animal{}, fmt.Errorf("animal not found")
		}
		return Animal{}, fmt.Errorf("error getting animal: %w", err)
	}

	return animal, nil
}

func (PostgresAnimalRepository) DeleteAnimal(id int) error {
	_, err := db.DB.Exec(`DELETE FROM animals WHERE id = $1`, id)

	return err
}
