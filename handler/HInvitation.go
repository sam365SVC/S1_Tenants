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
			"error": "invalid json format",
		})
	}
	if err := h._validator.Struct(&req); err != nil {
		errores := make(map[string]string)

		if validationsErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationsErrors {
				switch e.Tag() {
				case "required":
					errores[e.Field()] = "Este campo es obligatorio"

				case "email":
					errores[e.Field()] = "El formato del correo electrónico no es válido"

				case "work_context":
					// Este es el error de tu hook personalizado
					errores[e.Field()] = "La posición seleccionada no existe dentro del departamento indicado"

				case "gte":
					errores[e.Field()] = "El valor debe ser mayor o igual a " + e.Param()

				default:
					errores[e.Field()] = "Error de validación: " + e.Tag()
				}
			}
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "La validación de los datos ha fallado",
			"details": errores,
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
func (h *InvitationHandler) InvitationJob(c echo.Context) error {
	var req dtos.InvitationJobDto

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "invalid json format",
		})
	}

	tenantIdRaw := c.Get("tenant_id")
	tenantNameRaw := c.Get("tenant_name")
	if tenantIdRaw == nil || tenantNameRaw == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"success": false,
			"error":   "Missing authentication claims in token",
		})
	}
	tenantIdFloat, okTenantId := tenantIdRaw.(float64)
	tenantName, okTenantName := tenantNameRaw.(string)

	if !okTenantId || !okTenantName {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"success": false,
			"error":   "Invalid token claims format",
		})
	}
	t, status, err := h._service.SendInvitationJob(c.Request().Context(), tenantName, int(tenantIdFloat), req)
	if err != nil {
		c.Logger().Errorf("FALLO INTERNO - INVITATION JOB: %v|REQUEST: %v|HEADER ID:%d NAME:%s", err, int(tenantIdFloat), tenantName)
		switch status {
		case http.StatusConflict:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "Ya existe una invitación activa para este correo.",
			})

		case http.StatusServiceUnavailable:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "No pudimos enviar el correo de invitación. Por favor, intente más tarde.",
			})

		case http.StatusBadRequest:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "Los datos de la invitación son inválidos.",
			})

		case http.StatusInternalServerError:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "Ocurrió un error interno en el servidor.",
			})

		default:
			return c.JSON(status, echo.Map{
				"success": false,
				"error":   "Error inesperado, contacte con soporte técnico.",
			})
		}
	}
	return c.JSON(status, echo.Map{
		"success": true,
		"token":   t,
		"message": "invit sent with sucess!",
	})
}
