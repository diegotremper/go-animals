package animal

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var ErrAnimalNotFound = errors.New("animal not found")

type AnimalRepository interface {
	CreateAnimal(r AninalCreateRequest) error
	UpdateAnimal(id int64, r AnimalUpdateRequest) error
	ListAnimals() ([]Animal, error)
	GetAnimal(id int64) (Animal, error)
	DeleteAnimal(id int64) error
}

type PostgresAnimalRepository struct {
	log *zap.Logger
	db  *sqlx.DB
}

func NewPostgresAnimalRepository(log *zap.Logger, db *sqlx.DB) *PostgresAnimalRepository {
	return &PostgresAnimalRepository{log: log, db: db}
}

func (r *PostgresAnimalRepository) CreateAnimal(req AninalCreateRequest) error {
	_, err := r.db.Exec(`INSERT INTO animals (name, age, description) VALUES ($1, $2, $3)`, req.Name, req.Age, req.Description)
	if err != nil {
		return fmt.Errorf("failed to insert animal: %w", err)
	}
	r.log.Info("animal created", zap.String("name", req.Name), zap.Int("age", req.Age))
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
	r.log.Info("animal updated", zap.Int64("id", id))
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
	r.log.Info("listed animals", zap.Int("count", len(animals)))
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
	r.log.Info("retrieved animal", zap.Int64("id", animal.ID))
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
	r.log.Info("deleted animal", zap.Int64("id", id))
	return nil
}
