package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TenantId   int    `json:"tenant_id"`
// 	Departament string `json:"departament"`
// 	Position    string `json:"position"`
type InvitationEmployee struct {
	ent.Schema
}

func (InvitationEmployee)Fields()[]ent.Field{
	return []ent.Field{
		field.String("department").NotEmpty(),
		field.String("position").NotEmpty(),
	}
}

func (InvitationEmployee)Edges()[]ent.Edge{
	return []ent.Edge{
		edge.From("invitation",Invitation.Type).Ref("invitation_employee").Unique(),
		edge.From("tenant",Tenant.Type).Ref("invitation_employees").Unique(),
		edge.From("branch",Branch.Type).Ref("invitation_employees").Unique(),
	}
}
