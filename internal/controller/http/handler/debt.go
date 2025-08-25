package handler

import (
	"context"
	"log/slog"

	"shop/internal/entity"

	"github.com/gin-gonic/gin"
)

// DebtLogCreate godoc
// @Summary Create a new Debt Item
// @Description Create a new Debt Item with the provided details
// @Tags Debt
// @Accept  json
// @Produce  json
// @Param Debt body entity.DebtLogCreate true "Debt Details"
// @Success 200 {object} string
// @Failure 400 {object}  string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /debt/create [post]
func (h *Handler) DebtLogCreate(c *gin.Context) {
	reqBody := entity.DebtLogCreate{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.DebtLogsRepo.Create(context.Background(), &reqBody)
	if err != nil {
		c.JSON(500, gin.H{"Error creating Debt:": err})
		slog.Error("Error creating Debt: ", "err", err)
		return
	}

	slog.Info("New Debt created successfully")
	c.JSON(200, gin.H{"Massage": "New Debt created successfully"})
}

// GetDebts godoc
// @Summary Get Debt by ID
// @Description Get an Debt by their ID
// @Tags Debt
// @Accept  json
// @Produce  json
// @Param id query string false "User ID"
// @Param status query string false "Status"
// @Success 200 {object} entity.DebtLogGetAllRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /debt/get [get]
func (h *Handler) GetDebts(c *gin.Context) {
	user_id := c.Query("id")
	status := c.Query("status")

	res, err := h.UseCase.DebtLogsRepo.GetDebtLogs(context.Background(), user_id, status)
	if err != nil {
		c.JSON(500, gin.H{"Error getting Debt by ID: ": err})
		slog.Error("Error getting Debt by ID: ", "err", err)
		return
	}

	slog.Info("Debt retrieved successfully")
	c.JSON(200, res)
}

// UpdateDebt godoc
// @Summary Update an Debt
// @Description Update an Debt's details
// @Tags Debt
// @Accept  json
// @Produce  json
// @Param id query string true "Debt ID"
// @Param Debt body entity.DebtLogUpdateBody true "Debt Update Details"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /debt/update [put]
func (h *Handler) UpdateDebt(c *gin.Context) {
	reqBody := entity.DebtLogUpdateBody{}

	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body:": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.DebtLogsRepo.Update(context.Background(), entity.DebtLogUpdate{
		ID:     c.Query("id"),
		Amount: reqBody.Amount,
		Reason: reqBody.Reason,
		Status: reqBody.Status,
	})
	if err != nil {
		c.JSON(500, gin.H{"Error updating Debt:": err})
		slog.Error("Error updating Debt: ", "err", err)
		return
	}

	slog.Info("Debt updated successfully")
	c.JSON(200, "Debt updated successfully")
}

// Report godoc
// @Summary Generate Debt Report
// @Description Generate a report of Debts within a specified date range
// @Tags Debt
// @Accept  json
// @Produce  json
// @Param Debt body entity.Report true "Chek"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /debt/report [post]
func (h *Handler) Report(c *gin.Context) {
	reqBody := entity.Report{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		slog.Error("Error binding request body: ", "err", err)
		c.JSON(400, gin.H{"Error binding request body": err})
		return
	}

	err = h.UseCase.DebtLogsRepo.Report(context.Background(), &reqBody)
	if err != nil {
		c.JSON(500, gin.H{"Error generating report:": err})
		slog.Error("Error generating report: ", "err", err)
		return
	}
	slog.Info("Report generated successfully")
	c.JSON(200, "Report generated successfully")
}

// // DeleteDebtItem godoc
// // @Summary Delete an Debt
// // @Description Delete an Debt by ID
// // @Tags Debt
// // @Accept  json
// // @Produce  json
// // @Param id query string true "Debt Item ID"
// // @Success 200 {string} string "Debt deleted successfully"
// // @Failure 400 {object} string
// // @Failure 500 {object} string
// // @Security BearerAuth
// // @Router /debt/delete [delete]
// func (h *Handler) DeleteDebt(c *gin.Context) {
// 	Debt_id := c.Query("id")

// 	err := h.UseCase.DebtRepo.Delete(context.Background(), Debt_id)
// 	if err != nil {
// 		c.JSON(500, gin.H{"Error deleting Debt by ID:": err})
// 		slog.Error("Error deleting Debt by ID: ", "err", err)
// 		return
// 	}

// 	slog.Info("Debt deleted successfully")
// 	c.JSON(200, "Debt deleted successfully")
// }


