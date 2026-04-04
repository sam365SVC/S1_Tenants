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
		field.String("name").MaxLen(50).NotEmpty(),
		field.Time("end_suscription").SchemaType(map[string]string{
			dialect.Postgres:"date",
		}).Optional(),
	}
}
func (Tenant)Edges()[]ent.Edge  {
	return []ent.Edge{
		edge.From("plan",Plan.Type).Ref("tenants").Unique(),
		edge.To("employees",Employee.Type),
		edge.To("branchs",Branch.Type),
	}
}