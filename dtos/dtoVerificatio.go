package dtos

type VerificationDeveloperdto struct{
	Email	string	`json:"email" validate:"required,email"`
}