package animal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diegotremper/go-animals/internal/animal"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func mockListAnimals() ([]animal.Animal, error) {
	return []animal.Animal{
		{ID: 1, Name: "Cat", Age: 3, Description: "Domestic"},
		{ID: 2, Name: "Dog", Age: 5, Description: "Friendly"},
	}, nil
}

func TestListAnimalsHandler(t *testing.T) {
	origFunc := animal.ListAnimals
	animal.ListAnimals = mockListAnimals

	defer func() { animal.ListAnimals = origFunc }()

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	animal.ListAnimalsHandler(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Cat")
	assert.Contains(t, w.Body.String(), "Dog")
}
