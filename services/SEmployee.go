package services

import (
	"context"
	"fmt"
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/ent"
	"saas_identidad/ent/employee"
	"saas_identidad/ent/tenant"
	"saas_identidad/utils"
)

type EmployeeServices struct {
	_client *ent.Client
}

func NewEmployeeServices(client *ent.Client) *EmployeeServices {
	return &EmployeeServices{
		_client: client,
	}
}

func (s *EmployeeServices) GetEmployee(ctx context.Context, page int, pageSize int, tenantId int) (listEmployee []dtos.EmployeeResponseDto, status int, err error) {
	if page <= 0 {
		page = 1
	}
	query := s._client.Employee.Query()
	if tenantId > 0 {
		existTenant, err := s._client.Tenant.Query().Where(tenant.IDEQ(tenantId)).Exist(ctx)
		if err != nil {
			return []dtos.EmployeeResponseDto{}, http.StatusInternalServerError, fmt.Errorf("error searching tenant %d exist: %w", tenantId, err)
		}
		if !existTenant {
			return []dtos.EmployeeResponseDto{}, http.StatusNotFound, fmt.Errorf("the tenant with id: %d not exist", tenantId)
		}
		query = query.Where(employee.HasTenantWith(tenant.IDEQ(tenantId)))
	}
	countAll, err := query.Count(ctx)
	if err != nil {
		return []dtos.EmployeeResponseDto{}, http.StatusInternalServerError, fmt.Errorf("error count employee: %w", err)
	}
	if countAll == 0 {
		return []dtos.EmployeeResponseDto{}, http.StatusOK, nil
	}
	maxPage := utils.CalcularPagination(int64(countAll), pageSize)
	if page > maxPage {
		page = maxPage
	}
	offset := (page - 1) * pageSize
	list, err := query.
		Limit(pageSize).
		Offset(offset).
		WithTenant().
		WithEmail(func(eq *ent.EmailQuery) {
			eq.WithUser()
		}).
		All(ctx)
	if err != nil {
		return []dtos.EmployeeResponseDto{}, http.StatusInternalServerError, fmt.Errorf("error searching employees: %w", err)
	}
	for _, employee := range list {
		dateBirt := employee.Edges.Email.Edges.User.DateBirth.Format("02/01/2006")

		listEmployee = append(listEmployee, dtos.EmployeeResponseDto{
			Id:         employee.ID,
			TenantId:   employee.Edges.Tenant.ID,
			TenantName: employee.Edges.Tenant.Name,
			Department: employee.Department.String(),
			Posistion:  employee.Position,
			User: dtos.UserResponseDto{
				Id:        employee.Edges.Email.Edges.User.ID,
				Name:      employee.Edges.Email.Edges.User.Name,
				LastName:  employee.Edges.Email.Edges.User.LastName,
				CI:        employee.Edges.Email.Edges.User.Ci,
				DateBirth: dateBirt,
			},
		})
	}
	return listEmployee, http.StatusOK, nil
}
