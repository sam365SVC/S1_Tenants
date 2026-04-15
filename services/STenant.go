package services

import (
	"context"
	"fmt"
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/ent"
	"saas_identidad/utils"
)

type TenantService struct {
	_cliente *ent.Client
}

func NewTenantServices(client *ent.Client) *TenantService {
	return &TenantService{
		_cliente: client,
	}
}

func (s *TenantService) GetTenat(ctx context.Context, page int, pageSize int) (listTenant []dtos.TenantResponse, status int, err error) {
	if page < 1 {
		page = 1
	}
	coutAll, err := s._cliente.Tenant.Query().Count(ctx)
	if err != nil {
		return []dtos.TenantResponse{}, http.StatusInternalServerError, fmt.Errorf("error counting tenant: %w", err)
	}
	if coutAll == 0 {
		return []dtos.TenantResponse{}, http.StatusOK, nil
	}
	maxPage := utils.CalcularPagination(int64(coutAll), pageSize)
	if page > maxPage {
		page = maxPage
	}
	offset := (page - 1) * pageSize
	list, err := s._cliente.Tenant.Query().Limit(pageSize).Offset(offset).WithOwner().All(ctx)
	if err != nil {
		return []dtos.TenantResponse{}, http.StatusInternalServerError, fmt.Errorf("error get list tenant: %w", err)
	}
	for _, tenant := range list {
		var ownerId int
		var ownerName string
		var ownerLastName string

		if tenant.Edges.Owner != nil {
			ownerId = tenant.Edges.Owner.ID
			ownerName = tenant.Edges.Owner.Name
			ownerLastName = tenant.Edges.Owner.LastName
		}
		listTenant = append(listTenant, dtos.TenantResponse{
			Id:            tenant.ID,
			Name_Tenant:   tenant.Name,
			OwnerId:       ownerId,
			NameOwner:     ownerName,
			LastNameOwner: ownerLastName,
		})
	}
	return listTenant, http.StatusOK, nil
}
