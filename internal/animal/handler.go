package animal

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnimalHandler struct {
	repo AnimalRepository
}

func NewAnimalHandler(repo AnimalRepository) *AnimalHandler {
	return &AnimalHandler{repo: repo}
}

func (h *AnimalHandler) CreateAnimalHandler(ctx *gin.Context) {
	var req AninalCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	err := h.repo.CreateAnimal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create animal"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Animal created successfully"})
}

func (h *AnimalHandler) UpdateAnimalHandler(ctx *gin.Context) {
	var idStr = ctx.Param("id")
	var id, errA = strconv.Atoi(idStr)

	if errA != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var req AnimalUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	err := h.repo.UpdateAnimal(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to updating animal"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Animal updated successfully"})
}

func (h *AnimalHandler) ListAnimalsHandler(ctx *gin.Context) {
	var animals, err = h.repo.ListAnimals()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed retrieving list of animals"})
		return
	}

	ctx.JSON(http.StatusOK, animals)

}

func (h *AnimalHandler) GetAnimalHandler(ctx *gin.Context) {
	var idStr = ctx.Param("id")
	var id, errA = strconv.Atoi(idStr)
	if errA != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var animal, err = h.repo.GetAnimal(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting animal"})
		return
	}

	ctx.JSON(http.StatusOK, animal)
}

func (h *AnimalHandler) DeleteAnimalHandler(ctx *gin.Context) {
	var idStr = ctx.Param("id")
	var id, errA = strconv.Atoi(idStr)
	if errA != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err := h.repo.DeleteAnimal(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting animal"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Animal deleted successfully"})
}
