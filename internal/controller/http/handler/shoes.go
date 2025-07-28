package handler

import (
	"context"
	"log/slog"

	"shop/internal/entity"

	"github.com/gin-gonic/gin"
)

// CreateShoes godoc
// @Summary Create a new Shoes
// @Description Create a new Shoes with the provided details
// @Tags Shoes
// @Accept  json
// @Produce  json
// @Param Banner body entity.ShoesCreate true "Shoes Details"
// @Success 200 {object} entity.ShoesCreate
// @Failure 400 {object}  string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /shoes/create [post]
func (h *Handler) CreateShoes(c *gin.Context) {
	reqBody := entity.ShoesCreate{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.ShoesRepo.Create(context.Background(), &reqBody)
	if err != nil {
		c.JSON(500, gin.H{"Error creating Shoes:": err})
		slog.Error("Error creating Shoes: ", "err", err)
		return
	}

	slog.Info("New Shoes created successfully")
	c.JSON(200, gin.H{"Massage": "New Shoes created successfully"})
}

// GetByIdShoes godoc
// @Summary Get Shoes by ID
// @Description Get an Shoes by their ID
// @Tags Shoes
// @Accept  json
// @Produce  json
// @Param id query string true "Shoes ID"
// @Success 200 {object} entity.ShoesRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /shoes/get [get]
func (h *Handler) GetByIdShoes(c *gin.Context) {
	Shoes_id := c.Query("id")

	res, err := h.UseCase.ShoesRepo.GetById(context.Background(), &entity.ById{Id: Shoes_id})
	if err != nil {
		c.JSON(500, gin.H{"Error getting Shoes by ID: ": err})
		slog.Error("Error getting Shoes by ID: ", "err", err)
		return
	}

	slog.Info("Shoes retrieved successfully")
	c.JSON(200, res)
}

// UpdateShoes godoc
// @Summary Update an Shoes
// @Description Update an Shoes's details
// @Tags Shoes
// @Accept  json
// @Produce  json
// @Param id query string true "Shoes ID"
// @Param Shoes body entity.ShoesCreate true "Shoes Update Details"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /shoes/update [put]
func (h *Handler) UpdateShoes(c *gin.Context) {
	reqBody := entity.ShoesCreate{}

	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body:": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.ShoesRepo.Update(context.Background(), &entity.ShoesUpdate{
		Id:          c.Query("id"),
		Name:        reqBody.Name,
		Price:       reqBody.Price,
		Color:       reqBody.Color,
		ImgUrl:      reqBody.ImgUrl,
		Size:        reqBody.Size,
		CategoryId:  reqBody.CategoryId,
		Description: reqBody.Description,
	})
	if err != nil {
		c.JSON(500, gin.H{"Error updating Shoes:": err})
		slog.Error("Error updating Shoes: ", "err", err)
		return
	}

	slog.Info("Shoes updated successfully")
	c.JSON(200, "Shoes updated successfully")
}

// GetAllShoes godoc
// @Summary Get all Shoes
// @Description Get all Shoes with optional filtering
// @Tags Shoes
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} entity.ShoesGetAllRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /shoes/list [get]
func (h *Handler) GetAllShoes(c *gin.Context) {
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

	res, err := h.UseCase.ShoesRepo.GetAll(context.Background(), req)
	if err != nil {
		c.JSON(500, gin.H{"Error getting Shoes:": err})
		slog.Error("Error getting Shoes: ", "err", err)
		return
	}

	slog.Info("Shoes retrieved successfully")
	c.JSON(200, res)
}

// DeleteShoes godoc
// @Summary Delete an Shoes
// @Description Delete an Shoes by ID
// @Tags Shoes
// @Accept  json
// @Produce  json
// @Param id query string true "Shoes ID"
// @Success 200 {string} string "Shoes deleted successfully"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /shoes/delete [delete]
func (h *Handler) DeleteShoes(c *gin.Context) {
	Shoes_id := c.Query("id")

	err := h.UseCase.ShoesRepo.Delete(context.Background(), &entity.ById{Id: Shoes_id})
	if err != nil {
		c.JSON(500, gin.H{"Error deleting Shoes by ID:": err})
		slog.Error("Error deleting Shoes by ID: ", "err", err)
		return
	}

	slog.Info("Shoes deleted successfully")
	c.JSON(200, "Shoes deleted successfully")
}
