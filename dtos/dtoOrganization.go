package dtos

type OrganizationCreateDto struct{
	Tenant	TenantCreate	`json:"tenant"`
	Branch	BranchCreateDto	`json:"branch"`
}