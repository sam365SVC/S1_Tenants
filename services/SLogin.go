package services

import (
	"context"
	"fmt"
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/ent"
	"saas_identidad/ent/email"
	"saas_identidad/ent/employee"
	"saas_identidad/ent/tenant"
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

func (s *LoginServices) Login(ctx context.Context, req dtos.LoginDto) (Token string, list []*ent.Employee, status int, err error) {
	emailUser, err := s._client.Email.Query().Where(email.EmailEQ(req.Email)).WithUser().Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return "", nil, http.StatusUnauthorized, fmt.Errorf("credentials incorrect")
		}
		return "", nil, http.StatusInternalServerError, fmt.Errorf("error searching email: %w", err)
	}
	success := ChechPassword(req.Password, emailUser.PaswordHash)
	if !success {
		return "", nil, http.StatusUnauthorized, fmt.Errorf("credentials incorrect ")
	}
	if emailUser.Account.String() == "ADMIN" {
		employeesList, err := s._client.Employee.Query().
			Where(
				employee.HasEmailWith(email.ID(emailUser.ID)),
				employee.ActiveEQ(true)).
			WithTenant().
			All(ctx)
		if err != nil {
			return "", nil, http.StatusInternalServerError, fmt.Errorf("error searching employees: %w", err)
		}
		if len(employeesList) == 1 {
			tokenAdmin, err := security.GenerateTenantJWT(emailUser.Edges.User.ID, emailUser.Account.String(), emailUser.Email, employeesList[0].Edges.Tenant.ID, employeesList[0].Edges.Tenant.Name, employeesList[0].Department.String(), employeesList[0].Position)
			if err != nil {
				return "", nil, http.StatusInternalServerError, fmt.Errorf("error making tenant jwt: %w", err)
			}
			return tokenAdmin, nil, http.StatusOK, nil
		}
		if len(employeesList) > 1 {
			generalToken, err := security.GenerateGlobalJWT(emailUser.Edges.User.ID, emailUser.Account.String(), emailUser.Email)
			if err != nil {
				return "", nil, http.StatusInternalServerError, fmt.Errorf("error making global jwt: %w", err)
			}
			return generalToken, employeesList, http.StatusOK, nil
		}
	}
	generalToken, err := security.GenerateGlobalJWT(emailUser.Edges.User.ID, emailUser.Account.String(), emailUser.Email)
	if err != nil {
		return "", nil, http.StatusInternalServerError, fmt.Errorf("error making jwt: %w", err)
	}
	return generalToken, nil, http.StatusOK, nil
}
func (s *LoginServices) LoginTenant(ctx context.Context, userId int, userEmail string, req dtos.LoginTenantDto) (Token string, status int, err error) {
	employeeUser, err := s._client.Employee.Query().
		Where(
			employee.HasEmailWith(email.EmailEQ(userEmail)),
			employee.HasTenantWith(tenant.IDEQ(req.TargetTenantId)),
			employee.ActiveEQ(true),
		).
		WithTenant().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", http.StatusForbidden, fmt.Errorf("user does not have access to this tenant")
		}
		return "", http.StatusInternalServerError, fmt.Errorf("error verifying tenant access: %w", err)
	}
	newToken, err := security.GenerateTenantJWT(
		userId,
		"ADMIN",
		userEmail,
		employeeUser.Edges.Tenant.ID,
		employeeUser.Edges.Tenant.Name,
		employeeUser.Department.String(),
		employeeUser.Position,
	)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("error making tenant jwt: %w", err)
	}
	return newToken, http.StatusOK, nil
}
