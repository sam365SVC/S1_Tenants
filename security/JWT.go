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

	jwt.RegisteredClaims
}

func GenerateJWT(userId int, account, email string) (Token string,err error) {
	claims:=CustomClaimsJWT{
		UserId: userId,
		Account: account,
		Email: email,

		RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),
				IssuedAt: jwt.NewNumericDate(time.Now()),
				Issuer: "s1-tenant",
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	secretStr:=os.Getenv("JWT_KEY")
	if strings.TrimSpace(secretStr)=="" {
		return "",fmt.Errorf("error: don't have JWT_KEY")
	}

	jwtSecretStr:=[]byte(secretStr)

	tokenStr,err:=token.SignedString(jwtSecretStr)
	if err!=nil {
		return "",err
	}

	return tokenStr,nil
}