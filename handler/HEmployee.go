package handler

import (
	"net/http"
	"os"
	"saas_identidad/dtos"
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
func (h *EmployeeHandler) RemplaceEmployee(c echo.Context) error {
	employeeIdStr := c.Param("id")
	employeeId, err := strconv.Atoi(employeeIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "Invalid employee ID in the URL",
		})
	}
	var req dtos.EmployeeRemplaceDto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}
	if err := h._validator.Struct(&req); err != nil {
		errores := make(map[string]string)

		if validationsErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationsErrors {
				switch e.Tag() {
				case "required":
					errores[e.Field()] = "This field is mandatory"
				default:
					errores[e.Field()] = "Validation error: " + e.Tag()
				}
			}
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"success":  false,
			"error":    "Validation errors in the submitted data",
			"detalles": errores,
		})
	}
	res, status, err := h._service.RemplaceEmployee(c.Request().Context(), req, employeeId)

	if err != nil {
		c.Logger().Errorf("ERROR PUT EMPLOYEE: %v | ID: %d", err, employeeId)

		switch status {
		case http.StatusNotFound:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "The requested employee does not exist.",
			})
		case http.StatusBadRequest:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "Invalid data provided (e.g., invalid department enum).",
			})
		case http.StatusInternalServerError:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "An internal error occurred while replacing the employee data.",
			})
		default:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "An unexpected error occurred.",
			})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "Employee replaced successfully",
		"data":    res,
	})
}
func (h *EmployeeHandler) PatchEmployee(c echo.Context) error {
	employeeIdStr := c.Param("id")
	employeeId, err := strconv.Atoi(employeeIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "Invalid employee ID in the URL",
		})
	}

	var req dtos.EmployeePatchDto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "Invalid request format",
		})
	}

	if err := h._validator.Struct(&req); err != nil {
		errores := make(map[string]string)
		if validationsErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationsErrors {
				switch e.Tag() {
				default:
					errores[e.Field()] = "Validation error: " + e.Tag()
				}
			}
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success":  false,
			"error":    "Validation errors in the submitted data",
			"detalles": errores,
		})
	}

	res, status, err := h._service.PatchEmployee(c.Request().Context(), req, employeeId)

	if err != nil {
		c.Logger().Errorf("ERROR PATCH EMPLOYEE: %v | ID: %d", err, employeeId)

		switch status {
		case http.StatusNotFound:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "The requested employee does not exist.",
			})
		case http.StatusBadRequest:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "Invalid data provided (e.g., invalid department enum).",
			})
		case http.StatusInternalServerError:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "An internal error occurred while updating the employee data.",
			})
		default:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "An unexpected error occurred.",
			})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "Employee updated successfully",
		"data":    res,
	})
}
