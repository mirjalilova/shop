package handler

import (
	"context"
	"log/slog"

	"shop/internal/entity"

	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @Summary Create a new Category
// @Description Create a new Category with the provided details
// @Tags Category
// @Accept  json
// @Produce  json
// @Param Category body entity.CategoryCreate true "Category Details"
// @Success 200 {object} string
// @Failure 400 {object}  string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /category/create [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	reqBody := entity.CategoryCreate{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.CategoryRepo.Create(context.Background(), &reqBody)
	if err != nil {
		c.JSON(500, gin.H{"Error creating Category:": err})
		slog.Error("Error creating Category: ", "err", err)
		return
	}

	slog.Info("New Category created successfully")
	c.JSON(200, gin.H{"Massage": "New Category created successfully"})
}

// GetByIdCategory godoc
// @Summary Get Category by ID
// @Description Get an Category by their ID
// @Tags Category
// @Accept  json
// @Produce  json
// @Param id query string true "Category ID"
// @Success 200 {object} entity.CategoryRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /category/get [get]
func (h *Handler) GetByIdCategory(c *gin.Context) {
	Category_id := c.Query("id")

	res, err := h.UseCase.CategoryRepo.GetById(context.Background(), &entity.ById{Id: Category_id})
	if err != nil {
		c.JSON(500, gin.H{"Error getting Category by ID: ": err})
		slog.Error("Error getting Category by ID: ", "err", err)
		return
	}

	slog.Info("Category retrieved successfully")
	c.JSON(200, res)
}

// UpdateCategory godoc
// @Summary Update an Category
// @Description Update an Category's details
// @Tags Category
// @Accept  json
// @Produce  json
// @Param id query string true "Category ID"
// @Param Category body entity.CategoryCreate true "Category Update Details"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /category/update [put]
func (h *Handler) UpdateCategory(c *gin.Context) {
	reqBody := entity.CategoryUpdate{}

	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body:": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.CategoryRepo.Update(context.Background(), &entity.CategoryUpdate{
		Id:   c.Query("id"),
		Name: reqBody.Name,
	})
	if err != nil {
		c.JSON(500, gin.H{"Error updating Category:": err})
		slog.Error("Error updating Category: ", "err", err)
		return
	}

	slog.Info("Category updated successfully")
	c.JSON(200, "Category updated successfully")
}

// GetAllCategories godoc
// @Summary Get all Categories
// @Description Get all Categories with optional filtering
// @Tags Category
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} entity.CategoryGetAllRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /category/list [get]
func (h *Handler) GetAllCategories(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")

	limitValue, offsetValue, err := parsePaginationParams(c, limit, offset)
	if err != nil {
		c.JSON(400, gin.H{"Error parsing pagination parameters:": err.Error()})
		slog.Error("Error parsing pagination parameters: ", "err", err)
		return
	}

	req := &entity.Filter{
		Limit:  limitValue,
		Offset: offsetValue,
	}

	res, err := h.UseCase.CategoryRepo.GetAll(context.Background(), req)
	if err != nil {
		c.JSON(500, gin.H{"Error getting Categories:": err})
		slog.Error("Error getting Categories: ", "err", err)
		return
	}

	slog.Info("Categories retrieved successfully")
	c.JSON(200, res)
}

// DeleteCategory godoc
// @Summary Delete an Category
// @Description Delete an Category by ID
// @Tags Category
// @Accept  json
// @Produce  json
// @Param id query string true "Category ID"
// @Success 200 {string} string "Category deleted successfully"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /category/delete [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	Category_id := c.Query("id")

	err := h.UseCase.CategoryRepo.Delete(context.Background(), &entity.ById{Id: Category_id})
	if err != nil {
		c.JSON(500, gin.H{"Error deleting Category by ID:": err})
		slog.Error("Error deleting Category by ID: ", "err", err)
		return
	}

	slog.Info("Category deleted successfully")
	c.JSON(200, "Category deleted successfully")
}
