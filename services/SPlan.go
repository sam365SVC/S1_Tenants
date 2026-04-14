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
func (s *PlanServices) AllPlan(ctx context.Context)(list []*ent.Plan,status int,err error){
	listPlan,err:=s._client.Plan.Query().All(ctx)
	if err!=nil {
		return nil,http.StatusInternalServerError,fmt.Errorf("error get all plan: %w",err)
	}
	return listPlan,http.StatusOK,nil
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

func (s *PlanServices) UpdatePlan(ctx context.Context,req dtos.PlanUpdateDto,planId int)(plan *ent.Plan,status int, err error){
	updatePlan:=s._client.Plan.UpdateOneID(planId)

	if req.Subscription!=nil {
		updatePlan.SetSubscription(*req.Subscription)
	}
	if req.Price!=nil {
		updatePlan.SetPrice(*req.Price)
	}
	if req.MaxEmployees!=nil {
		updatePlan.SetMaxEmployees(*req.MaxEmployees)
	}
	if req.MaxBranches!=nil {
		updatePlan.SetMaxBranches(*req.MaxBranches)
	}
	if req.MaxBosses!=nil {
		updatePlan.SetMaxBosses(*req.MaxBosses)
	}

	planUpdate,err:=updatePlan.Save(ctx)
	if err!=nil {
		if ent.IsNotFound(err) {
			return nil,http.StatusNotFound,fmt.Errorf("error plan not fount: %w",err)
		}
		return nil,http.StatusInternalServerError,fmt.Errorf("error update plan: %w",err)
	}
	return planUpdate,http.StatusOK,nil
}
