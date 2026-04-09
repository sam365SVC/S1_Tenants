package services

import (
	"context"
	"fmt"
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/ent"
	"saas_identidad/ent/email"
	"saas_identidad/security"
)

type LoginServices struct {
	_client ent.Client
}

func NewLoginServices(client *ent.Client) *LoginServices {
	return &LoginServices{
		_client: *client,
	}
}

func (s *LoginServices) Login(ctx context.Context, req dtos.LoginDto) (Token string, status int, err error) {
	tx, err := s._client.Tx(ctx)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error making tx: %v", err)
	}
	emailUser, err := tx.Email.Query().Where(email.EmailEQ(req.Email)).Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return "", http.StatusUnauthorized, fmt.Errorf("credentials incorrect")
		}
		return "", http.StatusInternalServerError, fmt.Errorf("error searching email: %w", err)
	}
	user, err := emailUser.QueryUser().OnlyID(ctx)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error fetching user detalils: %w", err)
	}
	success := ChechPassword(req.Password, emailUser.PaswordHash)
	if !success {
		return "", http.StatusUnauthorized, fmt.Errorf("credentials incorrect ")
	}
	generalToken, err := security.GenerateJWT(user, emailUser.Account.String(), emailUser.Email)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error making jwt: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error making commit: %w", err)
	}
	return generalToken, http.StatusOK, nil
}
