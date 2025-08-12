package handler

import (
	"context"
	"log/slog"

	"shop/internal/entity"

	"github.com/gin-gonic/gin"
)

// BucketItemCreate godoc
// @Summary Create a new Bucket Item
// @Description Create a new Bucket Item with the provided details
// @Tags Bucket
// @Accept  json
// @Produce  json
// @Param Bucket body entity.BucketItemCreate true "Bucket Details"
// @Success 200 {object} string
// @Failure 400 {object}  string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /bucket/create [post]
func (h *Handler) BucketItemCreate(c *gin.Context) {
	reqBody := entity.BucketItemCreate{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.BucketRepo.Create(context.Background(), &reqBody)
	if err != nil {
		c.JSON(500, gin.H{"Error creating Bucket:": err})
		slog.Error("Error creating Bucket: ", "err", err)
		return
	}

	slog.Info("New Bucket created successfully")
	c.JSON(200, gin.H{"Massage": "New Bucket created successfully"})
}

// GetBucket godoc
// @Summary Get Bucket 
// @Description Get all Bucket 
// @Tags Bucket
// @Accept  json
// @Produce  json
// @Param id query string true "User ID"
// @Success 200 {object} entity.BucketRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /bucket/get [get]
func (h *Handler) GetBucket(c *gin.Context) {
	user_id := c.Query("id")

	res, err := h.UseCase.BucketRepo.GetBucket(context.Background(), user_id)
	if err != nil {
		c.JSON(500, gin.H{"Error getting Bucket by ID: ": err})
		slog.Error("Error getting Bucket by ID: ", "err", err)
		return
	}

	slog.Info("Bucket retrieved successfully")
	c.JSON(200, res)
}

// UpdateBucket godoc
// @Summary Update an Bucket
// @Description Update an Bucket's details
// @Tags Bucket
// @Accept  json
// @Produce  json
// @Param id query string true "Bucket ID"
// @Param Bucket body entity.BucketItemUpdateBody true "Bucket Update Details"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /bucket/update [put]
func (h *Handler) UpdateBucket(c *gin.Context) {
	reqBody := entity.BucketItemUpdateBody{}

	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body:": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.BucketRepo.Update(context.Background(), &entity.BucketItemUpdate{
		Id:    c.Query("id"),
		Count: reqBody.Count,
		Price: reqBody.Price,
	})
	if err != nil {
		c.JSON(500, gin.H{"Error updating Bucket:": err})
		slog.Error("Error updating Bucket: ", "err", err)
		return
	}

	slog.Info("Bucket updated successfully")
	c.JSON(200, "Bucket updated successfully")
}

// DeleteBucketItem godoc
// @Summary Delete an Bucket
// @Description Delete an Bucket by ID
// @Tags Bucket
// @Accept  json
// @Produce  json
// @Param id query string true "Bucket Item ID"
// @Success 200 {string} string "Bucket deleted successfully"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /bucket/delete [delete]
func (h *Handler) DeleteBucket(c *gin.Context) {
	bucket_id := c.Query("id")

	err := h.UseCase.BucketRepo.Delete(context.Background(), bucket_id)
	if err != nil {
		c.JSON(500, gin.H{"Error deleting Bucket by ID:": err})
		slog.Error("Error deleting Bucket by ID: ", "err", err)
		return
	}

	slog.Info("Bucket deleted successfully")
	c.JSON(200, "Bucket deleted successfully")
}
