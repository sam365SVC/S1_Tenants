package handler

import (
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	_service   *services.LoginServices
	_validator *validator.Validate
}

func NewAuthHandler(service *services.LoginServices, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		_service:   service,
		_validator: validator,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dtos.LoginDto

	// 1. Parsear el JSON
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid JSON format",
		})
	}

	// 2. Validar el Struct (Email y Password)
	if err := h._validator.Struct(&req); err != nil {
		errores := make(map[string]string)

		if validationsErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationsErrors {
				switch e.Tag() {
				case "required":
					errores[e.Field()] = "This field is mandatory"
				case "email":
					errores[e.Field()] = "Must be a valid email address"
				case "min":
					errores[e.Field()] = "Must have at least " + e.Param() + " characters"
				default:
					errores[e.Field()] = "Validation error: " + e.Tag()
				}
			}
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":   "Validation failed",
			"details": errores,
		})
	}

	// 3. Llamar al servicio de Login
	jwtT, status, err := h._service.Login(c.Request().Context(), req)

	if err != nil {
		// Log del error para el desarrollador (en consola)
		c.Logger().Errorf("INTERNAL ERROR - LOGIN: %v | Email: %s", err, req.Email)

		// Manejo de respuestas según el Status Code
		switch status {
		case http.StatusUnauthorized:
			// No decimos si falló el email o la clave por seguridad
			return c.JSON(status, echo.Map{
				"error": "Invalid email or password",
			})
		case http.StatusInternalServerError:
			return c.JSON(status, echo.Map{
				"error": "An internal server error occurred",
			})
		default:
			return c.JSON(status, echo.Map{
				"error": "Authentication failed, please contact support",
			})
		}
	}
	return c.JSON(status, echo.Map{
		"success": true,
		"token":   jwtT,
	})
}
