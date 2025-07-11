package animal

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnimalHandler struct {
	module Module
	repo   AnimalRepository
}

func NewAnimalHandler(module Module, repo AnimalRepository) *AnimalHandler {
	return &AnimalHandler{module: module, repo: repo}
}

func (h *AnimalHandler) CreateAnimalHandler(ctx *gin.Context) {
	var req AnimalCreateRequest
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
	idStr := ctx.Param("id")
	idInt, errA := strconv.Atoi(idStr)
	if errA != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	id := int64(idInt)

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
	idStr := ctx.Param("id")
	idInt, errA := strconv.Atoi(idStr)
	if errA != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	id := int64(idInt)

	var animal, err = h.repo.GetAnimal(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting animal"})
		return
	}

	ctx.JSON(http.StatusOK, animal)
}

func (h *AnimalHandler) DeleteAnimalHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idInt, errA := strconv.Atoi(idStr)
	if errA != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	id := int64(idInt)

	err := h.repo.DeleteAnimal(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting animal"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Animal deleted successfully"})
}
