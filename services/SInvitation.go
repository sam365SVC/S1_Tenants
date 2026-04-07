package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/ent"
	"saas_identidad/ent/email"
	"saas_identidad/ent/invitation"
	"saas_identidad/pkg/mailer"
)

type InvitationServices struct {
	_cliente *ent.Client
}

func NewInvitationServices(cliente *ent.Client) *InvitationServices{
	return &InvitationServices{
		_cliente: cliente,
	}
}
func generateToke()string  {
	b:=make([]byte,16)
	rand.Read(b)

	return hex.EncodeToString(b)
}

func (s *InvitationServices) VerificationDeveloper(ctx context.Context, req dtos.VerificationDeveloperdto, account string) (token string,status int, err error){
	tx,err:=s._cliente.Tx(ctx)
	if err!=nil {
		return "",http.StatusInternalServerError,fmt.Errorf("error to created tx: %w",err)
	}
	defer tx.Rollback()
	t:=generateToke()

	emailExist,err:=tx.Email.Query().Where(email.EmailEQ(req.Email)).Exist(ctx)
	if err!=nil {
		return "",http.StatusInternalServerError,fmt.Errorf("error consult if email exist in table email: %w",err)
	}
	if emailExist {
		return "",http.StatusConflict,fmt.Errorf("error the email is reguitre in db: %w",err)
	}

	err=mailer.ValidateEmailDeveloper(req.Email,t,account)
	if err!=nil {
		return "",http.StatusInternalServerError,fmt.Errorf("error to send email: %w",err)
	}
	_,err=tx.Invitation.Create().
		SetEmail(req.Email).
		SetToken(t).
		SetAccount(invitation.Account(account)).
		Save(ctx)
	if err!=nil {
		if ent.IsConstraintError(err) {
			return "",http.StatusConflict,fmt.Errorf("error invitation exist: %w",err)
		}
		if ent.IsValidationError(err) {
			return "",http.StatusBadRequest,fmt.Errorf("error request not valid: %w",err)
		}
		return "",http.StatusInternalServerError,fmt.Errorf("error to created invitation: %w",err)
	}
	if err:=tx.Commit();err!=nil {
		return "",http.StatusInternalServerError,fmt.Errorf("error when making commit: %w",err)
	}
	return t,http.StatusOK,nil
}