package handler

import (
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type PlanHandler struct {
	_service   *services.PlanServices
	_validator *validator.Validate
}

func NewPlanHandler(service *services.PlanServices, validator *validator.Validate) *PlanHandler {
	return &PlanHandler{
		_service:   service,
		_validator: validator,
	}
}

func (h PlanHandler) CreatePlan(c echo.Context) error {
	var req dtos.PlanCreateDto

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
				campoAfectado := e.Field()

				switch e.Tag() {
				case "max":
					errores[campoAfectado] = "It must have at most " + e.Param() + " characters"
				case "gte":
					errores[campoAfectado] = "The value must be greater than or equal to " + e.Param()
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
			"success": false,
			"error":   "Invalid request parameters",
		})
	}

	status, err := h._service.CreatePlan(c.Request().Context(), req)
	if err != nil {
		c.Logger().Errorf("FALLO INTERNO - PLAN CREATE: %v |Request: %v", err, req)
		switch status {
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
		"message": "plan create success",
	})
}
