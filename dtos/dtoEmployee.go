package dtos

type EmployeeResponseDto struct {
	Id         int             `json:"id"`
	TenantId   int             `json:"tenant_id"`
	TenantName string          `json:"tenant_name"`
	Department string          `json:"department"`
	Posistion  string          `json:"position"`
	User       UserResponseDto `json:"user,omitempty"`
}

type EmployeeRemplaceDto struct {
	Department string `json:"department" validate:"required"`
	Posistion  string `json:"position" validate:"required"`
	Active     bool   `json:"active" validate:"required"`
}
type EmployeePatchDto struct {
	Department string `json:"department" validate:"omitempty"`
	Posistion  string `json:"position" validate:"omitempty"`
	Active     *bool   `json:"active" validate:"omitempty"`
}
