package dtos

type LoginDto struct{
	Email	string	`json:"email" validate:"required,email"`
	Password	string	`json:"password" validate:"required,min=7"`
}

type LoginTenantDto struct{
	TargetTenantId	int	`json:"target_tenant_id" validate:"required,gt=0"`
}