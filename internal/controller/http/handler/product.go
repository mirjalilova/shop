package handler

import (
	"context"
	"log/slog"

	"shop/internal/entity"

	"github.com/gin-gonic/gin"
)

// CreateProduct godoc
// @Summary Create a new Product
// @Description Create a new Product with the provided details
// @Tags Product
// @Accept  json
// @Produce  json
// @Param Banner body entity.ProductCreate true "Product Details"
// @Success 200 {object} string
// @Failure 400 {object}  string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /product/create [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	reqBody := entity.ProductCreate{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.ProductRepo.Create(context.Background(), &reqBody)
	if err != nil {
		c.JSON(500, gin.H{"Error creating Product:": err})
		slog.Error("Error creating Product: ", "err", err)
		return
	}

	slog.Info("New Product created successfully")
	c.JSON(200, gin.H{"Massage": "New Product created successfully"})
}

// GetByIdProduct godoc
// @Summary Get Product by ID
// @Description Get an Product by their ID
// @Tags Product
// @Accept  json
// @Produce  json
// @Param id query string true "Product ID"
// @Success 200 {object} entity.ProductRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /product/get [get]
func (h *Handler) GetByIdProduct(c *gin.Context) {
	Product_id := c.Query("id")

	res, err := h.UseCase.ProductRepo.GetById(context.Background(), &entity.ById{Id: Product_id})
	if err != nil {
		c.JSON(500, gin.H{"Error getting Product by ID: ": err})
		slog.Error("Error getting Product by ID: ", "err", err)
		return
	}

	slog.Info("Product retrieved successfully")
	c.JSON(200, res)
}

// UpdateProduct godoc
// @Summary Update an Product
// @Description Update an Product's details
// @Tags Product
// @Accept  json
// @Produce  json
// @Param id query string true "Product ID"
// @Param Product body entity.ProductCreate true "Product Update Details"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /product/update [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	reqBody := entity.ProductCreate{}

	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body:": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.ProductRepo.Update(context.Background(), &entity.ProductUpdate{
		Id:          c.Query("id"),
		Name:        reqBody.Name,
		Price:       reqBody.Price,
		Type:        reqBody.Type,
		ImgUrl:      reqBody.ImgUrl,
		Size:        reqBody.Size,
		Count:       reqBody.Count,
		CategoryId:  reqBody.CategoryId,
		Description: reqBody.Description,
	})
	if err != nil {
		c.JSON(500, gin.H{"Error updating Product:": err})
		slog.Error("Error updating Product: ", "err", err)
		return
	}

	slog.Info("Product updated successfully")
	c.JSON(200, "Product updated successfully")
}

// GetAllProduct godoc
// @Summary Get all Product
// @Description Get all Product with optional filtering
// @Tags Product
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} entity.ProductGetAllRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /product/list [get]
func (h *Handler) GetAllProduct(c *gin.Context) {
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

	res, err := h.UseCase.ProductRepo.GetAll(context.Background(), req)
	if err != nil {
		c.JSON(500, gin.H{"Error getting Product:": err})
		slog.Error("Error getting Product: ", "err", err)
		return
	}

	slog.Info("Product retrieved successfully")
	c.JSON(200, res)
}

// DeleteProduct godoc
// @Summary Delete an Product
// @Description Delete an Product by ID
// @Tags Product
// @Accept  json
// @Produce  json
// @Param id query string true "Product ID"
// @Success 200 {string} string "Product deleted successfully"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /product/delete [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	Product_id := c.Query("id")

	err := h.UseCase.ProductRepo.Delete(context.Background(), &entity.ById{Id: Product_id})
	if err != nil {
		c.JSON(500, gin.H{"Error deleting Product by ID:": err})
		slog.Error("Error deleting Product by ID: ", "err", err)
		return
	}

	slog.Info("Product deleted successfully")
	c.JSON(200, "Product deleted successfully")
}
