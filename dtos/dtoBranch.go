package dtos

type BranchCreateDto struct{
	Name	string	`json:"name" validate:"required"`
	Phone 	string	`json:"phone" validate:"required"`
	Street	string	`json:"street" validate:"required,max=150"`
	Zone	string	`json:"zone" validate:"omitempty"`
	Reference	string	`json:"reference" validate:"omitempty,max=200"`
	City	string	`json:"city" validate:"required,max=50"`
	State	string	`json:"state" validate:"required,oneof=LP SC CB OR PT TJ CH BE PA"`
}

