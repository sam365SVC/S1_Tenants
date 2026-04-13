package handler

import (
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type InvitationHandler struct {
	_service   *services.InvitationServices
	_validator *validator.Validate
}

func NewInvitationHandler(service *services.InvitationServices, validator *validator.Validate) *InvitationHandler {
	return &InvitationHandler{
		_service:   service,
		_validator: validator,
	}
}
func (h *InvitationHandler) InvitationUserOrAdmin(c echo.Context) error {
	var req dtos.InvitationUserOrAdmindto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid json format",
		})
	}
	token, status, err := h._service.SendInvitation(c.Request().Context(), req.Email, req.Account)
	if err != nil {
		c.Logger().Errorf("FALLO INTERNO - INVITATION ADMIN OR USER: %v |Input: %v", err, req)

		switch status {
		case http.StatusConflict:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   err.Error(),
			})
		case http.StatusBadRequest:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "Invalid request data",
			})
		case http.StatusInternalServerError:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "internal server error",
			})
		default:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "contact with support now",
			})
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"token":   token,
	})
}

func (h *InvitationHandler) InvitationDeveloper(c echo.Context) error {
	var req dtos.InvitationDeveloper

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid json",
		})
	}
	token, status, err := h._service.SendInvitation(c.Request().Context(), req.Email, "DEVELOPER")
	if err != nil {
		c.Logger().Errorf("FALLO INTERNO - INVITATION DEVELOPER: %v|Input: %v", err, req)
		switch status {
		case http.StatusConflict:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   err.Error(),
			})
		case http.StatusBadRequest:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "Invalid request data",
			})
		case http.StatusInternalServerError:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "internal server error",
			})
		default:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "contact with support now",
			})
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"token":   token,
	})

}
