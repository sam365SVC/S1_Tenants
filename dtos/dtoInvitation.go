package dtos

type InvitationUserOrAdmindto struct {
	Email   string `json:"email" validate:"required,email"`
	Account string `json:"account" validate:"required,oneof=USER ADMIN"`
}

type InvitationDeveloper struct{
	Email string `json:"email" validate:"required,email"`	
}
