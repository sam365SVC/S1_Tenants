package dtos

type PlanCreateDto struct {
	Subscription string  `json:"subscription" validate:"max=20"`
	Price        float64 `json:"price" validate:"gte=0"`
	MaxEmployees int32   `json:"max_employees" validate:"gte=0"`
	MaxBranches  int32   `json:"max_branchs" validate:"gte=0"`
	MaxBosses    int32   `json:"max_bosses" validate:"gte=0"`
}

type PlanUpdateDto struct {
	Subscription *string  `json:"subscription,omitempty" validate:"omitempty,max=20"`
	Price        *float64 `json:"price,omitempty" validate:"omitempty,gte=0"`
	MaxEmployees *int32   `json:"max_employees,omitempty" validate:"omitempty,gte=0"`
	MaxBranches  *int32   `json:"max_branches,omitempty" validate:"omitempty,gte=0"` // Corregido a branches
	MaxBosses    *int32   `json:"max_bosses,omitempty" validate:"omitempty,gte=0"`
}
