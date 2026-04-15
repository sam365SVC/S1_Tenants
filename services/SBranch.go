package services

import (
	"context"
	"fmt"
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/ent"
	"saas_identidad/ent/branch"
	"saas_identidad/ent/tenant"
	"saas_identidad/utils"
)

type BranchServices struct {
	_client *ent.Client
}

func NewBranchServices(client *ent.Client) *BranchServices {
	return &BranchServices{
		_client: client,
	}
}

func (s *BranchServices) CreateBranch(ctx context.Context, req dtos.BranchCreateDto, ownerId int) (status int, err error) {
	tx, err := s._client.Tx(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error making tx: %w", err)
	}
	defer tx.Rollback()
	return http.StatusCreated, nil
}
func (s *BranchServices) GetPageBranch(ctx context.Context, page int, pageSize int, tenantId int) (listBranch []dtos.BranchResponseDto, status int, err error) {
	if page <= 0 {
		page = 1
	}
	query := s._client.Branch.Query()
	if tenantId > 0 {
		existTenant, err := s._client.Tenant.Query().Where(tenant.IDEQ(tenantId)).Exist(ctx)
		if err != nil {
			return []dtos.BranchResponseDto{}, http.StatusInternalServerError, fmt.Errorf("error searching tenant %d exist: %w", tenantId, err)
		}
		if !existTenant {
			return []dtos.BranchResponseDto{}, http.StatusNotFound, fmt.Errorf("the tenant with id: %d not exist", tenantId)
		}
		query = query.Where(branch.HasTenantWith(tenant.IDEQ(tenantId)))
	}
	countAll, err := query.Count(ctx)
	if err != nil {
		return []dtos.BranchResponseDto{}, http.StatusInternalServerError, fmt.Errorf("error count employee: %w", err)
	}
	if countAll == 0 {
		return []dtos.BranchResponseDto{}, http.StatusOK, nil
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
		All(ctx)
	if err != nil {
		return []dtos.BranchResponseDto{}, http.StatusInternalServerError, fmt.Errorf("error searching branch: %w", err)
	}
	for _, branch := range list {
		listBranch = append(listBranch, dtos.BranchResponseDto{
			Id:        branch.ID,
			Name:      branch.Name,
			Phone:     branch.Phone,
			Street:    branch.Street,
			Zone:      branch.Zone,
			Reference: branch.Reference,
			City:      branch.City,
			State:     branch.State.String(),
			Tenant: dtos.TenantResponse{
				Id:          branch.Edges.Tenant.ID,
				Name_Tenant: branch.Edges.Tenant.Name,
			},
		})
	}
	return listBranch, http.StatusOK, nil
}
