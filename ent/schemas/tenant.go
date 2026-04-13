package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Tenant struct {
	ent.Schema
}

func (Tenant) Fields()[]ent.Field{
	return []ent.Field{
		field.String("name").MaxLen(50).NotEmpty().Unique(),
		field.Time("end_suscription").SchemaType(map[string]string{
			dialect.Postgres:"date",
		}).Optional(),
	}
}
func (Tenant)Edges()[]ent.Edge  {
	return []ent.Edge{
		edge.From("owner",User.Type).Ref("organization").Unique().Required(),
		edge.From("plan",Plan.Type).Ref("tenants").Unique(),
		edge.To("invitation_employees",InvitationEmployee.Type),
		edge.To("employees",Employee.Type),
		edge.To("branches",Branch.Type),
	}
}