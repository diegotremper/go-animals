package animal

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateAnimalHandler(ctx *gin.Context) {
	var req AninalCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	err := CreateAnimal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create animal"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Animal created successfully"})
}

func UpdateAnimalHandler(ctx *gin.Context) {
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

	err := UpdateAnimal(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to updating animal"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Animal updated successfully"})
}

func ListAnimalsHandler(ctx *gin.Context) {
	var animals, err = ListAnimals()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed retrieving list of animals"})
		return
	}

	ctx.JSON(http.StatusOK, animals)

}

func GetAnimalHandler(ctx *gin.Context) {
	var idStr = ctx.Param("id")
	var id, errA = strconv.Atoi(idStr)
	if errA != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var animal, err = GetAnimal(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting animal"})
		return
	}

	ctx.JSON(http.StatusOK, animal)
}

func DeleteAnimalHandler(ctx *gin.Context) {
	var idStr = ctx.Param("id")
	var id, errA = strconv.Atoi(idStr)
	if errA != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err := DeleteAnimal(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting animal"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Animal deleted successfully"})
}
