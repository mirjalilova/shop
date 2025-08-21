package handler

import (
	"shop/internal/entity"

	"github.com/gin-gonic/gin"
)

// AddItem godoc
// @Summary Add an item to the Kassa
// @Description Add an item to the Kassa
// @Tags Kassa
// @Accept  json
// @Produce  json
// @Param id query string true "Product ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /kassa/add [post]
func (h *Handler) AddItem(c *gin.Context) {
	productID := c.Query("id")
	if productID == "" {
		c.JSON(400, gin.H{"error": "missing product ID"})
		return
	}

	err := h.UseCase.KassaRepo.AddItem(c.Request.Context(), productID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Item added successfully"})
}

// Formalize godoc
// @Summary Formalize an item in the Kassa
// @Description Formalize an item in the Kassa
// @Tags Kassa
// @Accept  json
// @Produce  json
// @Param id body []entity.Formalize true "Formalize body"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /kassa/formalize [post]
func (h *Handler) Formalize(c *gin.Context) {
	var formalize []entity.Formalize
	if err := c.ShouldBindJSON(&formalize); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	err := h.UseCase.KassaRepo.Formalize(c.Request.Context(), formalize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Items formalized successfully"})
}
