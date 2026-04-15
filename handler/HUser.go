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

type UserHandler struct {
	_service   *services.UserServices
	_validator *validator.Validate
}

func NewUserHandler(service *services.UserServices, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		_service:   service,
		_validator: validator,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dtos.UserCreateDto

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
					errores[e.Field()] = "This field is mandatory"
				case "email":
					errores[e.Field()] = "It must be a valid email"
				case "max":
					errores[e.Field()] = "It must have at most" + e.Param() + " characters"
				case "min":
					errores[e.Field()] = "It must have at least " + e.Param() + " characters"
				case "gt":
					errores[e.Field()] = "The value must be greater than " + e.Param()
				case "lte":
					errores[e.Field()] = "The value must be less than or equal to " + e.Param()
				case "datetime":
					// e.Param() nos dará el formato que falló (ej: 02/01/2006)
					errores[e.Field()] = "The date format is not valid (use the format:" + e.Param() + ")"

				// --- Validadores Personalizados ---
				case "is_name":
					errores[e.Field()] = "It must be a valid name (letters only and no special characters)"
				case "age_gte_16":
					errores[e.Field()] = "You must be at least 16 years old to register"
				case "secure_password":
					errores[e.Field()] = "The password is not secure enough"

				default:
					errores[e.Field()] = "Validation error: " + e.Tag()
				}
			}
		}

		// Retornamos el mapa de errores al frontend con un código 400 Bad Request
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":    "validation errors in the submitted data",
			"detalles": errores,
		})
	}
	status, err := h._service.CreateUser(c.Request().Context(), req)
	if err != nil {
		c.Logger().Errorf("FALLO INTERNO - USER CREATE: %v |Request: %v", err, req)
		switch status {
		case http.StatusUnauthorized:
			return c.JSON(status, echo.Map{
				"error": "token or email incorrect",
			})
		case http.StatusConflict:
			return c.JSON(status, echo.Map{
				"error": err,
			})
		case http.StatusBadRequest:
			return c.JSON(status, echo.Map{
				"error": err,
			})
		case http.StatusInternalServerError:
			return c.JSON(status, echo.Map{
				"error": "internal server error",
			})
		default:
			return c.JSON(status, echo.Map{
				"error": "contact with support now",
			})
		}
	}
	return c.JSON(status, echo.Map{
		"message": "user created with exit",
	})
}
func (h *UserHandler) AllUser(c echo.Context) error {
	pageStr := c.Param(":page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "invalid value param",
		})
	}
	sizeStr := os.Getenv("PageSize")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"error":   "error convert size",
		})
	}
	list, status, err := h._service.GetPageUser(c.Request().Context(), size, page)
	if err != nil {
		c.Logger().Errorf("FALLO INTERNO - USER CREATE: %v |Param: %w", err, page)
		switch status {
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
		"data":    list,
	})
}
func (h *UserHandler)GetUserId(c echo.Context)error{
	userIdStr:=c.Param("id")
	userId,err:=strconv.Atoi(userIdStr)
	if err!=nil {
		return c.JSON(http.StatusBadRequest,echo.Map{
			"success":false,
			"error":"error param is invalid",
		})
	}
	user,status,err:=h._service.GetUserId(c.Request().Context(),userId)
	if err!=nil {
		c.Logger().Errorf("FALLO INTERNO - USER CREATE: %v |Param: %s",err,userIdStr)
		switch status{
		case http.StatusBadRequest:
			return c.JSON(status,echo.Map{
				"success":false,
				"error": err,
			})
		case http.StatusNotFound:
			return c.JSON(status,echo.Map{
				"success":false,
				"error": "user %d not found",
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
	return c.JSON(status,echo.Map{
		"success":true,
		"user":user,
	})
}