package dtos

type TenantCreate struct {
	Name string `json:"name_corporation"`
}
type TenantUpdate struct {
}

type TenantResponse struct {
	Id            int    `json:"id"`
	Name_Tenant   string `json:"name_tenant"`
	OwnerId       int    `json:"owner_id,omitempty"`
	NameOwner     string `json:"owner_name,omitempty"`
	LastNameOwner string `json:"owner_last_name,omitempty"`
}
