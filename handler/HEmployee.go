package handler

import (
	"net/http"
	"os"
	"saas_identidad/services"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type EmployeeHandler struct {
	_service   *services.EmployeeServices
	_validator *validator.Validate
}

func NewEmployeeHandler(service *services.EmployeeServices, validator *validator.Validate) *EmployeeHandler {
	return &EmployeeHandler{
		_service:   service,
		_validator: validator,
	}
}

func (h *EmployeeHandler) GetPageEmployee(c echo.Context) error {
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "invalid page parameter",
		})
	}
	tenantIdStr := c.QueryParam("tenantId")
	var tenantId int
	if tenantIdStr != "" {
		tenantId, err = strconv.Atoi(tenantIdStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"success": false,
				"error":   "invalid tenantId query parameter",
			})
		}
	}
	sizeStr := os.Getenv("PageSize")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "error convert size",
		})
	}

	listEmployee, status, err := h._service.GetEmployee(c.Request().Context(), page, size, tenantId)
	if err != nil {
		c.Logger().Errorf("FALLO INTERNO - EMPLOYEE PAGE: %v |Param: %d|TenantId: %d", err, page, tenantId)
		switch status {
		case http.StatusNotFound:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "The requested tenant does not exist.",
			})

		case http.StatusInternalServerError:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "An internal error occurred while fetching the employees. Please try again later.",
			})

		default:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "An unexpected error occurred. Please contact support.",
			})
		}
	}
	return c.JSON(status, echo.Map{
		"success":   true,
		"message":   "list employee success",
		"employees": listEmployee,
	})
}
