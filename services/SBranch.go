package services

import (
	"context"
	"fmt"
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/ent"
)

type BranchServices struct {
	_client *ent.Client
}

func NewBranchServices(client *ent.Client) *BranchServices {
	return &BranchServices{
		_client: client,
	}
}

func (h *BranchServices) CreateBranch(ctx context.Context, req dtos.BranchCreateDto, ownerId int) (status int, err error) {
	tx, err := h._client.Tx(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error making tx: %w", err)
	}
	defer tx.Rollback()
	return http.StatusCreated, nil
}
