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
	emailUser, err := s._client.Email.Query().Where(email.EmailEQ(req.Email)).WithUser().Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return "", http.StatusUnauthorized, fmt.Errorf("credentials incorrect")
		}
		return "", http.StatusInternalServerError, fmt.Errorf("error searching email: %w", err)
	}
	success := ChechPassword(req.Password, emailUser.PaswordHash)
	if !success {
		return "", http.StatusUnauthorized, fmt.Errorf("credentials incorrect ")
	}
	generalToken, err := security.GenerateJWT(emailUser.Edges.User.ID, emailUser.Account.String(), emailUser.Email)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error making jwt: %w", err)
	}
	return generalToken, http.StatusOK, nil
}
