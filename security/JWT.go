package security

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaimsJWT struct {
	UserId  int    `json:"user_id"`
	Account string `json:"account"`
	Email   string `json:"email"`

	TenantId   int    `json:"tenant_id,omitempty"`
	TenantName string `json:"tenant_name,omitempty"`
	Department string `json:"department,omitempty"`
	Position   string `json:"position,omitempty"`

	jwt.RegisteredClaims
}

func GenerateGlobalJWT(userId int, account, email string) (Token string, err error) {
	claims := CustomClaimsJWT{
		UserId:  userId,
		Account: account,
		Email:   email,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "s1-tenant",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretStr := os.Getenv("JWT_KEY")
	if strings.TrimSpace(secretStr) == "" {
		return "", fmt.Errorf("error: don't have JWT_KEY")
	}

	jwtSecretStr := []byte(secretStr)

	tokenStr, err := token.SignedString(jwtSecretStr)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func GenerateTenantJWT(userId int, account, email string, tenant_id int, tenant_name, department, position string) (Token string, err error) {
	claims := CustomClaimsJWT{
		UserId:     userId,
		Account:    account,
		Email:      email,
		TenantId:   tenant_id,
		TenantName: tenant_name,
		Department: department,
		Position:   position,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "s1-tenant",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretStr := os.Getenv("JWT_KEY")
	if strings.TrimSpace(secretStr) == "" {
		return "", fmt.Errorf("error: don't have JWT_KEY")
	}

	jwtSecretStr := []byte(secretStr)

	tokenStr, err := token.SignedString(jwtSecretStr)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
