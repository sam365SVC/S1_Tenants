package handler

import (
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type InvitationHandler struct {
	_service *services.InvitationServices
	_validator *validator.Validate
}

func NewInvitationHandler(service *services.InvitationServices, validator *validator.Validate) *InvitationHandler{
	return &InvitationHandler{
		_service: service,
		_validator: validator,
	}
}
func (h *InvitationHandler) InvitationDeveloper(c echo.Context) error{
	var req dtos.VerificationDeveloperdto
	account:="DEVELOPER"
	if err:=c.Bind(&req);err!=nil {
		return c.JSON(http.StatusBadRequest,echo.Map{
			"error":"invalid json format",
		})
	}
	token,status,err:=h._service.VerificationDeveloper(c.Request().Context(),req,account)
	if err!=nil {
		c.Logger().Errorf("FALLO INTERNO - INVITATION DEVELOPER OR USER: %v |Input: %v",err,req)

		switch status{
		case http.StatusConflict:
			return c.JSON(status,echo.Map{
				"error":"email exist in api",
			})
		case http.StatusBadRequest:
			return c.JSON(status,echo.Map{
				"error":"error request not valid",
			})
		case http.StatusInternalServerError:
			return c.JSON(status,echo.Map{
				"error":"internal server error",
			})
		default:
			return c.JSON(status,echo.Map{
				"error":"contact with support now",
			})
		}
	}
	return c.JSON(http.StatusOK,echo.Map{
		"ok":"email send",
		"token":token,
	})
}