package services

import (
	"context"
	"fmt"
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/ent"
	"saas_identidad/ent/email"
	"saas_identidad/ent/invitation"
	"saas_identidad/ent/user"
	"time"
)

type UserServices struct {
	_client *ent.Client
}

func NewUserServices(client *ent.Client) *UserServices {
	return &UserServices{
		_client: client,
	}
}

func (s *UserServices) CreateUser(ctx context.Context, req dtos.UserCreateDto) (status int, err error) {
	tx, err := s._client.Tx(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error creating transaction: %w", err)
	}
	defer tx.Rollback()
	invit, err := tx.Invitation.Query().Where(invitation.TokenEQ(req.Token)).Only(ctx)
	if err != nil {
		return http.StatusUnauthorized, fmt.Errorf("error: token is invalid or already used: %w", err)
	}
	if invit.Email != req.Email {
		return http.StatusUnauthorized, fmt.Errorf("error: email not autorized")
	}
	if invit.ExpireAt.Before(time.Now()) {
		return http.StatusUnauthorized, fmt.Errorf("error: token has expired")
	}
	datesUser, err := tx.User.Query().Where(user.CiEQ(req.CI)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return http.StatusInternalServerError, fmt.Errorf("error searching for user: %w", err)
	}
	//make password_hash
	password_hash, err := HashPassword(req.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if datesUser != nil {
		hasSameAccountType, err := datesUser.QueryEmails().
			Where(email.AccountEQ(email.Account(invit.Account))).Exist(ctx)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error querying user accounts: %w",err)
		}
		if hasSameAccountType {
			return http.StatusConflict, fmt.Errorf("error the user with CI:%d have an account: %s", datesUser.Ci, invit.Account)
		}
	}else{
		if req.Name==""||req.LastName==""||req.DateBirth==""{
			return http.StatusBadRequest,fmt.Errorf("error the user don't exist in the api(name, last_name and date_birth are require)")
		}
		dateBirth, err := time.Parse("02/01/2006", req.DateBirth)
		if err!=nil {
			return http.StatusBadRequest, fmt.Errorf("invalid date format, use DD/MM/YYYY: %w", err)
		}
		datesUser, err = tx.User.Create().
			SetName(req.Name).
			SetLastName(req.LastName).
			SetCi(req.CI).
			SetDateBirth(dateBirth).
			Save(ctx)
		if err != nil {
			if ent.IsConstraintError(err) {
				return http.StatusConflict, fmt.Errorf("Error the user exist: %w", err)
			}
			return http.StatusInternalServerError, fmt.Errorf("error to create user: %w", err)
		}
	}	
	_, err = tx.Email.Create().
		SetEmail(req.Email).
		SetPaswordHash(password_hash).
		SetAccount(email.Account(invit.Account)).
		SetUser(datesUser).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return http.StatusConflict, fmt.Errorf("error email exist in DB: %w", err)
		}
		if ent.IsValidationError(err) {
			return http.StatusBadRequest, err
		}
		return http.StatusInternalServerError, fmt.Errorf("error to created email: %w", err)
	}
	if err:=tx.Invitation.DeleteOne(invit).Exec(ctx);err!=nil {
		return http.StatusInternalServerError,fmt.Errorf("error when delete invitation: %w",err)
	}
	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error when making commit: %w", err)
	}
	return http.StatusCreated, nil
}