package security

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "missing autization header",
				})
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "invalid or expired token",
				})
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				c.Set("user_id", claims["user_id"])
				c.Set("account", claims["account"])
				c.Set("email", claims["email"])
				// los datos que faltan
				c.Set("tenant_id", claims["tenant_id"])
				c.Set("tenant_name", claims["tenant_name"])
				c.Set("department", claims["department"])
				c.Set("position", claims["position"])
			}
			return next(c)
		}
	}
}

func RequireRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userAccount, ok := c.Get("account").(string)
			if !ok || userAccount == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "could not identify the user's role (invalid token)",
				})
			}
			isAllowed := false

			for _, role := range allowedRoles {
				if userAccount == role {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "could not identify the user's role (invalid token)",
				})
			}
			return next(c)
		}
	}
}

// RequireOfficeBoss verifica que el usuario pertenezca al departamento 'office' y sea 'boss'
func RequireOfficeBoss() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userToken := c.Get("user").(*jwt.Token)

			claims, ok := userToken.Claims.(*CustomClaimsJWT)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "error al parsear claims")
			}
			// Usamos strings.ToLower para evitar problemas de mayúsculas/minúsculas
			isOffice := strings.ToLower(claims.Department) == "office"
			isBoss := strings.ToLower(claims.Position) == "boss"

			if !isOffice || !isBoss {
				return echo.NewHTTPError(http.StatusForbidden, "Acceso denegado: se requiere ser Boss de Office")
			}
			return next(c)
		}
	}
}
