package handler

import (
	"net/http"
	"os"
	"saas_identidad/services"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type BranchHandler struct {
	_service   *services.BranchServices
	_validator *validator.Validate
}

func NewBranchHandler(service *services.BranchServices, validator *validator.Validate) *BranchHandler {
	return &BranchHandler{
		_service:   service,
		_validator: validator,
	}
}

func (h *BranchHandler) GetPageBranch(c echo.Context) error {
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "Invalid page parameter. It must be a number.",
		})
	}

	tenantIdStr := c.QueryParam("tenantId")
	var tenantId int
	if tenantIdStr != "" {
		tenantId, err = strconv.Atoi(tenantIdStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"success": false,
				"error":   "Invalid tenantId parameter. It must be a number.",
			})
		}
	}
	pageSizeStr := os.Getenv("PageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	list, status, err := h._service.GetPageBranch(c.Request().Context(), page, pageSize, tenantId)

	if err != nil {
		c.Logger().Errorf("ERROR GET BRANCH PAGE: %v | tenantId: %d", err, tenantId)

		switch status {
		case http.StatusNotFound:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "The requested tenant does not exist.",
			})
		case http.StatusInternalServerError:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "An internal error occurred while fetching the branches. Please try again later.",
			})
		default:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "An unexpected error occurred. Please contact support.",
			})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"data":    list,
	})
}
