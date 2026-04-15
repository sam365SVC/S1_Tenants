package dtos

type EmployeeResponseDto struct {
	Id         int             `json:"id"`
	TenantId   int             `json:"tenant_id"`
	TenantName string          `json:"tenant_name"`
	Department string          `json:"department"`
	Posistion  string          `json:"position"`
	User       UserResponseDto `json:"user"`
}
