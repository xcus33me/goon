package user

import (
	"goon/service/auth"
	"goon/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) handleLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "login",
	})
}

func (h *Handler) handleRegister(c *gin.Context) {
	var payload types.RegisterUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	if payload.Name == "" || payload.Email == "" || payload.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "All fields are required",
		})
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check user existence",
		})
		return
	}

	if user != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
	}

	hashedPassword, err := auth.HashPassword(payload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
	}

	verificationToken := uuid.New().String()
	tokenExpiresAt := time.Now().Add(25 * time.Hour)

	newUser := types.User{
		ID:                uuid.New(),
		Name:              payload.Name,
		Email:             payload.Email,
		Password:          hashedPassword,
		VerificationToken: verificationToken,
		TokenExpiresAt:    tokenExpiresAt,
		Role:              types.RoleUser,
	}

	if err := h.store.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": newUser.ID,
	})
}
