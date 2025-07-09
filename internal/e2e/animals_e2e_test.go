package e2e_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/diegotremper/go-animals/infrastructure"
	"github.com/diegotremper/go-animals/internal/animal"

	"github.com/docker/go-connections/nat"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	mpostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func startTestServer(db *sql.DB) (*http.Server, string, func(context.Context) error, error) {
	gin.SetMode(gin.TestMode)
	r := infrastructure.SetupRouter(infrastructure.InitLogger(), sqlx.NewDb(db, "postgres"))
	// Use dynamic port
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to listen: %w", err)
	}
	addr := ln.Addr().String()
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go srv.Serve(ln)
	// baseURL for http requests
	baseURL := fmt.Sprintf("http://%s", addr)
	return srv, baseURL, srv.Shutdown, nil
}

func applyMigrations(db *sql.DB, t *testing.T) {
	driver, err := mpostgres.WithInstance(db, &mpostgres.Config{})
	if err != nil {
		t.Fatalf("failed to create postgres driver: %v", err)
	}
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to determine project root: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/migrations", projectRoot),
		"postgres", driver)
	if err != nil {
		t.Fatalf("failed to create migrator: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("migration failed: %v", err)
	}
}

func createAnimal(t *testing.T, baseURL string) {
	resp, err := http.Post(baseURL+"/animals", "application/json",
		strings.NewReader(`{"name":"E2ETest","age":5,"description":"e2e check"}`))
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func listAnimals(t *testing.T, baseURL string) {
	resp, err := http.Get(baseURL + "/animals")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var animals []animal.Animal
	if err := json.Unmarshal(body, &animals); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(animals) == 0 {
		t.Fatal("expected at least one animal")
	}
	found := false
	for _, a := range animals {
		if a.Name == "E2ETest" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("E2ETest animal not found in list")
	}
}

func getAnimal(t *testing.T, db *sql.DB, baseURL string) int {
	var id int
	row := db.QueryRow("SELECT id FROM animals WHERE name = $1", "E2ETest")
	if err := row.Scan(&id); err != nil {
		t.Fatalf("failed to fetch animal ID: %v", err)
	}
	resp, err := http.Get(fmt.Sprintf("%s/animals/%d", baseURL, id))
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var a animal.Animal
	if err := json.Unmarshal(body, &a); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if a.ID != id || a.Name != "E2ETest" {
		t.Fatalf("unexpected animal data: %+v", a)
	}
	return id
}

func deleteAnimal(t *testing.T, baseURL string, id int) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/animals/%d", baseURL, id), nil)
	if err != nil {
		t.Fatalf("create delete req failed: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on delete, got %d", resp.StatusCode)
	}
}

func assertAnimalDeleted(t *testing.T, db *sql.DB, id int) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM animals WHERE id = $1", id)
	if err := row.Scan(&count); err != nil {
		t.Fatalf("failed count check: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 records, got %d", count)
	}
}

func TestE2E_AnimalsLifecycle(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:15.2",
		postgres.WithDatabase("animals"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForSQL("5432", "postgres", func(host string, port nat.Port) string {
				return fmt.Sprintf("host=%s port=%s user=test password=test dbname=animals sslmode=disable", host, port.Port())
			}).WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("could not start container: %v", err)
	}
	defer pgContainer.Terminate(ctx)

	connStr := pgContainer.MustConnectionString(ctx, "sslmode=disable")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	applyMigrations(db, t)

	_, baseURL, shutdown, err := startTestServer(db)
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer shutdown(ctx)
	time.Sleep(2 * time.Second)

	t.Run("create animal", func(t *testing.T) {
		createAnimal(t, baseURL)
	})
	t.Run("list animals and check count", func(t *testing.T) {
		listAnimals(t, baseURL)
	})
	var createdID int
	t.Run("get animal by ID", func(t *testing.T) {
		createdID = getAnimal(t, db, baseURL)
	})
	t.Run("delete animal", func(t *testing.T) {
		deleteAnimal(t, baseURL, createdID)
	})
	t.Run("ensure animal is gone", func(t *testing.T) {
		assertAnimalDeleted(t, db, createdID)
	})
}
