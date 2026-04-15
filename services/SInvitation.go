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
func (s *InvitationServices) SendInvitationJob(ctx context.Context, tenant_name string, tenantId int, req dtos.InvitationJobDto) (token string, status int, err error) {
	tx, err := s._cliente.Tx(ctx)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error making tx: %w", err)
	}
	// El defer rollback es vital por si algo falla a mitad de camino
	defer tx.Rollback()

	// 1. Verificar si el USUARIO ya tiene cuenta en el sistema global
	// Importante: Aquí buscas en la tabla 'User', no en 'Invitation'
	userExists, err := tx.Email.Query().Where(email.EmailEQ(req.Email)).Exist(ctx)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error checking user table: %w", err)
	}

	// 2. Verificar si ya hay una INVITACIÓN pendiente para este correo
	invitExist, err := tx.Invitation.Query().Where(invitation.EmailEQ(req.Email)).Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return "", http.StatusInternalServerError, fmt.Errorf("error searching for invitation: %w", err)
		}
	} else {
		// Si la invitación existe y no ha expirado, evitamos duplicados
		if invitExist.ExpireAt.After(time.Now().UTC()) {
			return "", http.StatusConflict, fmt.Errorf("a valid invitation already exists")
		}
		// Si expiró, la borramos
		if err := tx.Invitation.DeleteOne(invitExist).Exec(ctx); err != nil {
			return "", http.StatusInternalServerError, fmt.Errorf("error deleting expired invitation: %w", err)
		}
	}

	// 3. Crear la nueva invitación
	t := generateToke() // Asegúrate de que esta función esté definida
	invitationCreate, err := tx.Invitation.Create().
		SetEmail(req.Email).
		SetToken(t).
		SetAccount("ADMIN"). // Generalmente las invitaciones de empleo son rol USER
		Save(ctx)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error creating invitation: %w", err)
	}

	// 4. Crear los detalles del empleo (Relación con Tenant y Branch)
	_, err = tx.InvitationEmployee.Create().
		SetDepartment(req.Department).
		SetPosition(req.Position).
		SetInvitation(invitationCreate).
		SetTenantID(tenantId).
		SetBranchID(req.Branch_id).
		Save(ctx)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error making invitation employee: %w", err)
	}

	// 5. ENVIAR EL CORREO (Lógica de bifurcación)
	// Lo hacemos antes del Commit si queremos que la DB falle si el mail falla,
	// o después si preferimos asegurar la DB primero.
	if userExists {
		// CASO 1: Ya es usuario del sistema. Usamos la plantilla de "Añadido a organización"
		err = mailer.SendInviteToOrganization(t, req.Email, tenant_name, req.Department, req.Position)
	} else {
		// CASO 2: Es un usuario nuevo. Usamos la plantilla de "Crear cuenta"
		err = mailer.ValidateJob(req.Email, tenant_name, t, req.Department, req.Position)
	}

	if err != nil {
		// Si el correo falla, cancelamos la transacción para no dejar datos huérfanos
		return "", http.StatusServiceUnavailable, fmt.Errorf("failed to send email: %w", err)
	}

	// 6. Consolidar cambios
	if err := tx.Commit(); err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error making commit: %w", err)
	}
	return t, http.StatusCreated, nil
}
