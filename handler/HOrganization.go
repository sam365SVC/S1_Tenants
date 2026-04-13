package handler

import (
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type OrganizationHandler struct {
	_service   *services.OrganizationServices
	_validator *validator.Validate
}

func NewOrganizationHandler(service *services.OrganizationServices, validator *validator.Validate) *OrganizationHandler {
	return &OrganizationHandler{
		_service:   service,
		_validator: validator,
	}
}

func (h *OrganizationHandler) CreateOrganization(c echo.Context) error {
	var req dtos.OrganizationCreateDto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "invalid data format",
		})
	}
	if err := h._validator.Struct(&req); err != nil {
		errores := make(map[string]string)

		if validationsErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationsErrors {
				campoAfectado := e.StructNamespace()

				switch e.Tag() {
				case "required":
					errores[campoAfectado] = "This field is mandatory"
				case "max":
					errores[campoAfectado] = "It must have at most " + e.Param() + " characters"
				case "oneof":
					errores[campoAfectado] = "The value must be one of: " + e.Param()
				default:
					errores[campoAfectado] = "Validation error: " + e.Tag()
				}
			}
			return c.JSON(http.StatusBadRequest, echo.Map{
				"success": false,
				"error":   "Validation failed",
				"details": errores,
			})
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"sucess": false,
			"error":  "Invalid request parameters",
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

	if !okEmail || !okId {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"success": false,
			"error":   "Invalid token claims format",
		})
	}
	status, err := h._service.CreateOrganization(c.Request().Context(), int(userIdFloat), email, req)
	if err != nil {
		c.Logger().Errorf("FALLO INTERNO - ORGANIZATION CREATE: %v |Request: %v", err, req)
		switch status {
		case http.StatusConflict:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "The resource already exists (conflict with organization or branch name).",
			})

		case http.StatusNotFound:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "The provided email does not exist in our records.",
			})

		case http.StatusInternalServerError:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "Internal server error. Please try again later.",
			})

		default:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "An unexpected error occurred. Please contact support now.",
			})
		}
	}
	return c.JSON(status, echo.Map{
		"success": true,
		"message": "organization create success",
	})
}
