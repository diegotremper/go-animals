package animal_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/diegotremper/go-animals/internal/animal"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct{}

func (m mockRepo) ListAnimals() ([]animal.Animal, error) {
	return []animal.Animal{
		{ID: 1, Name: "Cat", Age: 3, Description: "Domestic"},
		{ID: 2, Name: "Dog", Age: 5, Description: "Friendly"},
	}, nil
}

func (m mockRepo) CreateAnimal(r animal.AninalCreateRequest) error           { return nil }
func (m mockRepo) UpdateAnimal(id int64, r animal.AnimalUpdateRequest) error { return nil }
func (m mockRepo) GetAnimal(id int64) (animal.Animal, error) {
	return animal.Animal{ID: id, Name: "Lion", Age: 7, Description: "Fierce"}, nil
}
func (m mockRepo) DeleteAnimal(id int64) error { return nil }

func (m mockRepo) CreateAnimalFail(r animal.AninalCreateRequest) error {
	return errors.New("failed to create")
}
func (m mockRepo) UpdateAnimalFail(id int64, r animal.AnimalUpdateRequest) error {
	return errors.New("failed to update")
}
func (m mockRepo) GetAnimalFail(id int64) (animal.Animal, error) {
	return animal.Animal{}, errors.New("not found")
}
func (m mockRepo) DeleteAnimalFail(id int64) error { return errors.New("failed to delete") }

type mockFailRepo struct {
	mockRepo
}

func (m mockFailRepo) CreateAnimal(r animal.AninalCreateRequest) error {
	return m.CreateAnimalFail(r)
}
func (m mockFailRepo) UpdateAnimal(id int64, r animal.AnimalUpdateRequest) error {
	return m.UpdateAnimalFail(id, r)
}
func (m mockFailRepo) GetAnimal(id int64) (animal.Animal, error) {
	return m.GetAnimalFail(id)
}
func (m mockFailRepo) DeleteAnimal(id int64) error {
	return m.DeleteAnimalFail(id)
}

func TestListAnimalsHandler(t *testing.T) {
	repo := mockRepo{}
	handler := animal.NewAnimalHandler(repo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	handler.ListAnimalsHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Cat")
	assert.Contains(t, w.Body.String(), "Dog")
}

func TestCreateAnimalHandler(t *testing.T) {
	repo := mockRepo{}
	handler := animal.NewAnimalHandler(repo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("POST", "/animals", strings.NewReader(`{"name":"Tiger","age":4,"description":"Wild"}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.CreateAnimalHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Animal created successfully")
}

func TestUpdateAnimalHandler(t *testing.T) {
	repo := mockRepo{}
	handler := animal.NewAnimalHandler(repo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	ctx.Request = httptest.NewRequest("PUT", "/animals/1", strings.NewReader(`{"name":"Panther","age":6,"description":"Stealthy"}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateAnimalHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Animal updated successfully")
}

func TestGetAnimalHandler(t *testing.T) {
	repo := mockRepo{}
	handler := animal.NewAnimalHandler(repo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.GetAnimalHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Lion")
}

func TestDeleteAnimalHandler(t *testing.T) {
	repo := mockRepo{}
	handler := animal.NewAnimalHandler(repo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	handler.DeleteAnimalHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Animal deleted successfully")
}

func TestCreateAnimalHandler_Failure(t *testing.T) {
	repo := mockFailRepo{}
	handler := animal.NewAnimalHandler(repo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest("POST", "/animals", strings.NewReader(`{"name":"Tiger","age":4,"description":"Wild"}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.CreateAnimalHandler(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create animal")
}

func TestUpdateAnimalHandler_Failure(t *testing.T) {
	repo := mockFailRepo{}
	handler := animal.NewAnimalHandler(repo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	ctx.Request = httptest.NewRequest("PUT", "/animals/1", strings.NewReader(`{"name":"Panther","age":6,"description":"Stealthy"}`))
	ctx.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateAnimalHandler(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to updating animal")
}

func TestGetAnimalHandler_Failure(t *testing.T) {
	repo := mockFailRepo{}
	handler := animal.NewAnimalHandler(repo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "99"}}

	handler.GetAnimalHandler(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error getting animal")
}

func TestDeleteAnimalHandler_Failure(t *testing.T) {
	repo := mockFailRepo{}
	handler := animal.NewAnimalHandler(repo)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Params = gin.Params{gin.Param{Key: "id", Value: "99"}}

	handler.DeleteAnimalHandler(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error deleting animal")
}
