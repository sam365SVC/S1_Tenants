package services

import (
	"context"
	"fmt"
	"net/http"
	"saas_identidad/dtos"
	"saas_identidad/ent"
)

type PlanServices struct {
	_client ent.Client
}

func NewPlanServices(client *ent.Client) *PlanServices {
	return &PlanServices{
		_client: *client,
	}
}

func (s *PlanServices) CreatePlan(ctx context.Context, req dtos.PlanCreateDto) (status int, err error) {
	tx, err := s._client.Tx(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error making tx: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.Plan.Create().
		SetSubscription(req.Subscription).
		SetPrice(req.Price).
		SetMaxEmployees(req.MaxEmployees).
		SetMaxBranches(req.MaxBranches).
		SetMaxBosses(req.MaxBosses).
		Save(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error making plan: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error making commit: %w", err)
	}
	return http.StatusCreated, nil
}
