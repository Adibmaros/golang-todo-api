package handler

import (
	"errors"
	"net/http"

	"github.com/Adibmaros/todo-api/internal/dto"
	"github.com/Adibmaros/todo-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Register godoc
// @Summary      Registrasi user baru
// @Description  Mendaftarkan akun user baru dengan email dan password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterRequest  true  "Data registrasi"
// @Success      201      {object}  model.User
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body tidak valid"})
		return
	}

	user, err := h.service.Register(input.Email, input.Password)
	if errors.Is(err, service.ErrEmailAlreadyUsed) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mendaftarkan user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary      Login user
// @Description  Autentikasi user dan mengembalikan JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Kredensial login"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request body tidak valid"})
		return
	}

	token, err := h.service.Login(input.Email, input.Password)
	if errors.Is(err, service.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
