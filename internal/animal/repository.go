package animal

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var ErrAnimalNotFound = errors.New("animal not found")

type AnimalRepository interface {
	CreateAnimal(r AnimalCreateRequest) error
	UpdateAnimal(id int64, r AnimalUpdateRequest) error
	ListAnimals() ([]Animal, error)
	GetAnimal(id int64) (Animal, error)
	DeleteAnimal(id int64) error
}

type PostgresAnimalRepository struct {
	db *sqlx.DB
}

func NewPostgresAnimalRepository(db *sqlx.DB) *PostgresAnimalRepository {
	return &PostgresAnimalRepository{db: db}
}

func (r *PostgresAnimalRepository) CreateAnimal(req AnimalCreateRequest) error {
	_, err := r.db.Exec(`INSERT INTO animals (name, age, description) VALUES ($1, $2, $3)`, req.Name, req.Age, req.Description)
	if err != nil {
		return fmt.Errorf("failed to insert animal: %w", err)
	}
	return nil
}

func (r *PostgresAnimalRepository) UpdateAnimal(id int64, req AnimalUpdateRequest) error {
	res, err := r.db.Exec(`UPDATE animals SET name = $1, age = $2, description = $3 WHERE id = $4`, req.Name, req.Age, req.Description, id)
	if err != nil {
		return fmt.Errorf("failed to update animal: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected on update: %w", err)
	}
	if rows == 0 {
		return ErrAnimalNotFound
	}
	return nil
}

func (r *PostgresAnimalRepository) ListAnimals() ([]Animal, error) {
	var (
		animals      []Animal = make([]Animal, 0)
		sqlStatement          = `SELECT id, name, age, description FROM animals`
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

func (r *PostgresAnimalRepository) GetAnimal(id int64) (Animal, error) {
	var (
		animal       Animal
		sqlStatement = `SELECT * FROM animals WHERE id = $1`
	)

	err := r.db.QueryRowx(sqlStatement, id).StructScan(&animal)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Animal{}, fmt.Errorf("%w: id=%d", ErrAnimalNotFound, id)
		}
		return Animal{}, fmt.Errorf("error getting animal: %w", err)
	}
	return animal, nil
}

func (r *PostgresAnimalRepository) DeleteAnimal(id int64) error {
	res, err := r.db.Exec(`DELETE FROM animals WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete animal: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected on delete: %w", err)
	}
	if rows == 0 {
		return ErrAnimalNotFound
	}
	return nil
}
