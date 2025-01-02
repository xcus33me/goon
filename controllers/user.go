package controllers

import (
	repo "goon/repository"
	"goon/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {
	var userInput struct {
		Name     string
		Email    string
		Password string
	}

	if err := c.ShouldBindJSON(&userInput); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"validations": utils.FormatValidationErrors(errs),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !utils.IsUniqueValue("users", "email", userInput.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"validations": map[string]interface{}{
				"Email": "The email is already exist!",
			},
		})
		return
	}

	verificationToken := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)

	hashPassword, err := utils.HashPassword(userInput.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user, err := repo.SaveUser(userInput.Name, userInput.Email, hashPassword, verificationToken.String(), expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
