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
			"success": false,
			"error":   "Invalid JSON format",
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
	jwtT, list, status, err := h._service.Login(c.Request().Context(), req)

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
	if list != nil {
		return c.JSON(status, echo.Map{
			"success":   true,
			"token":     jwtT,
			"employees": list,
		})
	}
	return c.JSON(status, echo.Map{
		"success": true,
		"token":   jwtT,
	})
}
func (h *AuthHandler) LoginTenant(c echo.Context) error {
	var req dtos.LoginTenantDto

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "Invalid JSON format",
		})
	}
	if err := h._validator.Struct(&req); err != nil {
		errores := make(map[string]string)

		if validationsErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationsErrors {
				switch e.Tag() {
				case "required":
					errores[e.Field()] = "This field is mandatory"
				case "gte":
					errores[e.Field()] = "The value must be greater than or equal to " + e.Param()
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
	userIdRaw := c.Get("user_id")
	emailRaw := c.Get("email")

	if userIdRaw == nil || emailRaw == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"success": false,
			"error":   "Missing authentication claims in token",
		})
	}
	userIdFloat, okId := userIdRaw.(float64)
	email, okEmail := emailRaw.(string)

	if !okId || !okEmail {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"success": false,
			"error":   "Invalid token claims format",
		})
	}
	tokenEmployee, status, err := h._service.LoginTenant(c.Request().Context(), int(userIdFloat), email, req)
	if err != nil {
		c.Logger().Errorf("FALLO INTERNO - LOGIN TENANT: %v |Request: %v", err, req)
		switch status {
		case http.StatusForbidden:
			// Caso: ent.IsNotFound -> El usuario no existe en esa empresa o está inactivo
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "You do not have access to this organization or your employment is inactive.",
			})

		case http.StatusInternalServerError:
			// Caso: Falló la base de datos o falló la generación del JWT
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "Internal server error. Please try again later.",
			})

		default:
			// Fallback de seguridad por si en el futuro agregas más códigos en el servicio
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "An unexpected error occurred. Please contact support now.",
			})
		}
	}
	return c.JSON(status, echo.Map{
		"success": true,
		"message": "login tenant success",
		"token":   tokenEmployee,
	})
}
