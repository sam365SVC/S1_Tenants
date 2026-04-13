package services

import (
	"context"
	"fmt"
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/ent"
	"saas_identidad/ent/branch"
	"saas_identidad/ent/email"
	"saas_identidad/ent/employee"
)

type OrganizationServices struct {
	_client *ent.Client
}

func NewOrganizationServices(client *ent.Client) *OrganizationServices {
	return &OrganizationServices{
		_client: client,
	}
}

func (s *OrganizationServices) CreateOrganization(ctx context.Context, userId int, emailR string, req dtos.OrganizationCreateDto) (status int, err error) {
	tx, err := s._client.Tx(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error making tx: %w", err)
	}
	defer tx.Rollback()
	tenantCreated, err := tx.Tenant.Create().
		SetName(req.Tenant.Name).
		SetOwnerID(userId).
		SetPlanID(1).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return http.StatusConflict, fmt.Errorf("error conflict in tenant: %w", err)
		}
		return http.StatusInternalServerError, fmt.Errorf("error making tenant: %w", err)
	}
	branchCreated, err := tx.Branch.Create().
		SetName(req.Branch.Name).
		SetPhone(req.Branch.Phone).
		SetStreet(req.Branch.Street).
		SetZone(req.Branch.Zone).
		SetReference(req.Branch.Reference).
		SetCity(req.Branch.City).
		SetState(branch.State(req.Branch.State)).
		SetTenant(tenantCreated).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return http.StatusConflict, fmt.Errorf("error conflict in branch: %w", err)
		}
		return http.StatusInternalServerError, fmt.Errorf("error making branch: %w", err)
	}
	emailEntity, err := tx.Email.Query().Where(email.EmailEQ(emailR)).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return http.StatusNotFound, fmt.Errorf("error email not exist: %w", err)
		}
		return http.StatusInternalServerError, fmt.Errorf("error found email: %w", err)
	}

	_, err = tx.Employee.Create().
		SetDepartment(employee.Department("office")).
		SetPosition("boss").
		SetEmails(emailEntity).
		SetTenant(tenantCreated).
		SetBranches(branchCreated).
		Save(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error making employee: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error making commit: %w", err)
	}

	return http.StatusCreated, nil
}
