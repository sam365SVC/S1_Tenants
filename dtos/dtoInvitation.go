package dtos

type InvitationUserOrAdmindto struct {
	Email   string `json:"email" validate:"required,email"`
	Account string `json:"account" validate:"required,oneof=USER ADMIN"`
}

type InvitationDeveloper struct{
	Email string `json:"email" validate:"required,email"`	
}

type InvitationJobDto struct{
	Branch_id	int		`json:"branch_id" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Department string	`json:"department" validate:"required"`
	Position 	string	`json:"position" validate:"required,work_context"`
}
