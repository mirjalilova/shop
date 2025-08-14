package handler

import (
	"context"
	"log/slog"
	"strconv"

	"shop/internal/entity"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body entity.CreateUser true "User Details"
// @Success 200 {object} string
// @Failure 400 {object}  string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /users/create [post]
func (h *Handler) Register(c *gin.Context) {
	reqBody := entity.CreateUser{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.UserRepo.Create(context.Background(), &reqBody)
	if err != nil {
		c.JSON(500, gin.H{"Error creating user:": err})
		slog.Error("Error creating user: ", "err", err)
		return
	}

	slog.Info("New user created successfully")
	c.JSON(200, gin.H{"Massage": "User registered successfully"})
}

// Login godoc
// @Summary User login
// @Description User login with phone number and password
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body entity.LoginReq true "User Login Details"
// @Success 200 {object} entity.LoginRes
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /users/login [post]
func (h *Handler) Login(c *gin.Context) {
	reqBody := entity.LoginReq{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	res, err := h.UseCase.UserRepo.Login(context.Background(), &reqBody)
	if err != nil {
		c.JSON(500, gin.H{"Error logging in user:": err})
		slog.Error("Error logging in user: ", "err", err)
		return
	}

	slog.Info("User logged in successfully")
	c.JSON(200, res)
}

// GetByIdUser godoc
// @Summary Get User by ID
// @Description Get a User by their ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id query string true "User ID"
// @Success 200 {object} entity.UserInfo
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /users/get [get]
func (h *Handler) GetByIdUser(c *gin.Context) {
	User_id := c.Query("id")

	res, err := h.UseCase.UserRepo.GetById(context.Background(), &entity.ById{Id: User_id})
	if err != nil {
		c.JSON(500, gin.H{"Error getting User by ID: ": err})
		slog.Error("Error getting User by ID: ", "err", err)
		return
	}

	slog.Info("User retrieved successfully")
	c.JSON(200, res)
}

// UpdateUser godoc
// @Summary Update a User
// @Description Update a User's details
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id query string true "User ID"
// @Param User body entity.UpdateUserBody true "User Update Details"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /users/update [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	reqBody := entity.UpdateUserBody{}

	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{"Error binding request body:": err})
		slog.Error("Error binding request body: ", "err", err)
		return
	}

	err = h.UseCase.UserRepo.Update(context.Background(), &entity.UpdateUser{
		Id:          c.Query("id"),
		Name:        reqBody.Name,
		PhoneNumber: reqBody.PhoneNumber,
	})
	if err != nil {
		c.JSON(500, gin.H{"Error updating User:": err})
		slog.Error("Error updating User: ", "err", err)
		return
	}

	slog.Info("User updated successfully")
	c.JSON(200, "User updated successfully")
}

// GetAllUsers godoc
// @Summary Get all Users
// @Description Get all Users with optional filtering
// @Tags Users
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param name query int false "Search by name"
// @Success 200 {object} entity.UserList
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /users/list [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
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

	res, err := h.UseCase.UserRepo.GetAll(context.Background(), req, c.Query("name"))
	if err != nil {
		c.JSON(500, gin.H{"Error getting Users:": err})
		slog.Error("Error getting Users: ", "err", err)
		return
	}

	slog.Info("Users retrieved successfully")
	c.JSON(200, res)
}

// DeleteUser godoc
// @Summary Delete a User
// @Description Delete a User by ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id query string true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /users/delete [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	User_id := c.Query("id")

	err := h.UseCase.UserRepo.Delete(context.Background(), &entity.ById{Id: User_id})
	if err != nil {
		c.JSON(500, gin.H{"Error deleting User by ID:": err})
		slog.Error("Error deleting User by ID: ", "err", err)
		return
	}

	slog.Info("User deleted successfully")
	c.JSON(200, "User deleted successfully")
}

func parsePaginationParams(c *gin.Context, limit, offset string) (int, int, error) {
	limitValue := 10
	offsetValue := 0

	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil {
			slog.Error("Invalid limit value", "err", err)
			c.JSON(400, gin.H{"error": "Invalid limit value"})
			return 0, 0, err
		}
		limitValue = parsedLimit
	} else {
		limitValue = 0
	}

	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil {
			slog.Error("Invalid offset value", "err", err)
			c.JSON(400, gin.H{"error": "Invalid offset value"})
			return 0, 0, err
		}
		offsetValue = parsedOffset
	}

	return limitValue, offsetValue, nil
}
