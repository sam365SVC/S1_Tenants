package handler

import (
	"net/http"
	"os"
	"saas_identidad/services"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type TenantHandler struct {
	_service   *services.TenantService
	_validator *validator.Validate
}

func NewTenantHandler(service *services.TenantService, validator *validator.Validate) *TenantHandler {
	return &TenantHandler{
		_service:   service,
		_validator: validator,
	}
}

func (h *TenantHandler) GetPageTenant(c echo.Context) error {
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "param is invalid",
		})
	}
	pageSizeStr := os.Getenv("PageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "invalid entorn variable",
		})
	}

	list, status, err := h._service.GetTenat(c.Request().Context(), page, pageSize)
	if err != nil {
		c.Logger().Errorf("FALLO INTERNO - TENANT GET: %v |Param: %s", err, pageStr)
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
		"message": "pagination success",
		"tenants": list,
	})
}
