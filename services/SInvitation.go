package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"saas_identidad/ent"
	"saas_identidad/ent/email"
	"saas_identidad/ent/invitation"
	"saas_identidad/pkg/mailer"
	"time"
)

type InvitationServices struct {
	_cliente *ent.Client
}

func NewInvitationServices(cliente *ent.Client) *InvitationServices {
	return &InvitationServices{
		_cliente: cliente,
	}
}
func generateToke() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *InvitationServices) SendInvitation(ctx context.Context, emailSend, account string) (token string, status int, err error) {
	tx, err := s._cliente.Tx(ctx)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error to created tx: %w", err)
	}
	defer tx.Rollback()
	t := generateToke()

	emailExist, err := tx.Email.Query().Where(email.EmailEQ(emailSend)).Exist(ctx)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error consult email: %w", err)
	}
	if emailExist {
		return "", http.StatusConflict, fmt.Errorf("error the email is reguitre in db: %w", err)
	}

	invitExist, err := tx.Invitation.Query().Where(invitation.EmailEQ(emailSend)).Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return "", http.StatusInternalServerError, fmt.Errorf("error searching for user: %w", err)
		}
	} else {
		if invitExist.ExpireAt.After(time.Now().UTC()) {
			return "", http.StatusConflict, fmt.Errorf("a valid invitation already exists and has not expired")
		}

		// Si llegó aquí, existe pero ya expiró. La borramos para crear una nueva.
		if err := tx.Invitation.DeleteOne(invitExist).Exec(ctx); err != nil {
			return "", http.StatusInternalServerError, fmt.Errorf("error deleting expired invitation: %w", err)
		}
	}
	_, err = tx.Invitation.Create().
		SetEmail(emailSend).
		SetToken(t).
		SetAccount(invitation.Account(account)).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return "", http.StatusConflict, fmt.Errorf("error invitation exist: %w", err)
		}
		if ent.IsValidationError(err) {
			return "", http.StatusBadRequest, fmt.Errorf("error request not valid: %w", err)
		}
		return "", http.StatusInternalServerError, fmt.Errorf("error to created invitation: %w", err)
	}
	go func() {
		err = mailer.ValidateEmailDeveloper(emailSend, t, account)
	}()

	if err := tx.Commit(); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error when making commit: %w", err)
	}
	return t, http.StatusOK, nil
}
