package animal

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AnimalRepository interface {
	CreateAnimal(r AninalCreateRequest) error
	UpdateAnimal(id int, r AnimalUpdateRequest) error
	ListAnimals() ([]Animal, error)
	GetAnimal(id int) (Animal, error)
	DeleteAnimal(id int) error
}

type PostgresAnimalRepository struct {
	db *sqlx.DB
}

func NewPostgresAnimalRepository(db *sqlx.DB) *PostgresAnimalRepository {
	return &PostgresAnimalRepository{db: db}
}

func (r *PostgresAnimalRepository) CreateAnimal(req AninalCreateRequest) error {
	_, err := r.db.Exec(`INSERT INTO animals (name, age, description) VALUES ($1, $2, $3)`, req.Name, req.Age, req.Description)

	return err
}

func (r *PostgresAnimalRepository) UpdateAnimal(id int, req AnimalUpdateRequest) error {
	_, err := r.db.Exec(`UPDATE animals SET name = $1, age = $2, description = $3 WHERE id = $4`, req.Name, req.Age, req.Description, id)

	return err
}

func (r *PostgresAnimalRepository) ListAnimals() ([]Animal, error) {
	var (
		animals      []Animal = make([]Animal, 0)
		sqlStatement          = `
			SELECT id, name, age, description FROM animals
		`
	)

	rows, err := r.db.Queryx(sqlStatement)
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

func (r *PostgresAnimalRepository) GetAnimal(id int) (Animal, error) {
	var (
		animal       Animal
		sqlStatement = `SELECT * FROM animals WHERE id = $1`
	)

	err := r.db.QueryRowx(sqlStatement, id).StructScan(&animal)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Animal{}, fmt.Errorf("animal not found")
		}
		return Animal{}, fmt.Errorf("error getting animal: %w", err)
	}

	return animal, nil
}

func (r *PostgresAnimalRepository) DeleteAnimal(id int) error {
	_, err := r.db.Exec(`DELETE FROM animals WHERE id = $1`, id)

	return err
}
