package handler

import (
	"context"
	"log/slog"

	"shop/internal/entity"

	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
// @Summary Create a new Order
// @Description Create a new Order with the provided details
// @Tags Order
// @Accept  json
// @Produce  json
// @Param Order body entity.OrderCreate true "Order Details"
// @Success 200 {object} string
// @Failure 400 {object}  string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /order/create [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	reqBody := entity.OrderCreate{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.OrderRepo.Create(context.Background(), &reqBody)
	if err != nil {
		c.JSON(500, gin.H{"Error creating Order:": err})
		slog.Error("Error creating Order: ", "err", err)
		return
	}

	slog.Info("New Order created successfully")
	c.JSON(200, gin.H{"Massage": "New Order created successfully"})
}

// GetOrders godoc
// @Summary Get Orders
// @Description Get all Orders
// @Tags Order
// @Accept  json
// @Produce  json
// @Param id query string true "User ID"
// @Param status query string true "Status"
// @Success 200 {object} entity.OrderRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /order/get [get]
func (h *Handler) GetOrders(c *gin.Context) {
	user_id := c.Query("id")

	res, err := h.UseCase.OrderRepo.GetOrders(context.Background(), user_id, c.Query("status"))
	if err != nil {
		c.JSON(500, gin.H{"Error getting Order by ID: ": err})
		slog.Error("Error getting Order by ID: ", "err", err)
		return
	}

	slog.Info("Order retrieved successfully")
	c.JSON(200, res)
}
