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
func (s *EmployeeServices) RemplaceEmployee(ctx context.Context, req dtos.EmployeeRemplaceDto, employeeId int) (res dtos.EmployeeResponseDto, status int, err error) {
	employeeRemplace, err := s._client.Employee.Query().
		Where(employee.IDEQ(employeeId)).
		WithTenant().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return dtos.EmployeeResponseDto{}, http.StatusNotFound, fmt.Errorf("employee not exist: %w", err)
		}
		return dtos.EmployeeResponseDto{}, http.StatusInternalServerError, fmt.Errorf("error searching employee: %w", err)
	}

	employeeUpdate, err := employeeRemplace.Update().
		SetDepartment(employee.Department(req.Department)).
		SetPosition(req.Posistion).
		SetActive(req.Active).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return dtos.EmployeeResponseDto{}, http.StatusBadRequest, fmt.Errorf("invalid constraint data for employee: %w", err)
		}
		return dtos.EmployeeResponseDto{}, http.StatusInternalServerError, fmt.Errorf("error updating employee: %w", err)
	}

	res = dtos.EmployeeResponseDto{
		Id:         employeeUpdate.ID,
		TenantId:   employeeRemplace.Edges.Tenant.ID,
		TenantName: employeeRemplace.Edges.Tenant.Name,
		Department: employeeUpdate.Department.String(),
		Posistion:  employeeUpdate.Position,
	}
	return res, http.StatusOK, nil
}

func (s *EmployeeServices) PatchEmployee(ctx context.Context, req dtos.EmployeePatchDto, employeeId int) (res dtos.EmployeeResponseDto, status int, err error) {
	
	employeeBase, err := s._client.Employee.Query().
		Where(employee.IDEQ(employeeId)).
		WithTenant().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return dtos.EmployeeResponseDto{}, http.StatusNotFound, fmt.Errorf("employee not exist: %w", err)
		}
		return dtos.EmployeeResponseDto{}, http.StatusInternalServerError, fmt.Errorf("error searching employee: %w", err)
	}

	updater := employeeBase.Update()
	if req.Department != "" {
		updater.SetDepartment(employee.Department(req.Department))
	}

	if req.Posistion != "" {
		updater.SetPosition(req.Posistion)
	}

	if req.Active != nil {
		updater.SetActive(*req.Active) // Desreferenciamos con * para obtener el valor true/false
	}

	employeeUpdate, err := updater.Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return dtos.EmployeeResponseDto{}, http.StatusBadRequest, fmt.Errorf("invalid constraint data for employee (e.g. invalid department): %w", err)
		}
		return dtos.EmployeeResponseDto{}, http.StatusInternalServerError, fmt.Errorf("error updating employee: %w", err)
	}

	res = dtos.EmployeeResponseDto{
		Id:         employeeUpdate.ID,
		TenantId:   employeeBase.Edges.Tenant.ID,
		TenantName: employeeBase.Edges.Tenant.Name,
		Department: employeeUpdate.Department.String(),
		Posistion:  employeeUpdate.Position,
	}

	return res, http.StatusOK, nil
}
