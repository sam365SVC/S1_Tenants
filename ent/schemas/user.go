package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(50).
			Validate(ValidateName).
			NotEmpty(),
		field.String("last_name").MaxLen(50).
			Validate(ValidateName).
			NotEmpty(),
		field.Int("ci").Range(11111, 999999999).Unique(),
		field.String("phone").MaxLen(14).Optional(),
		field.Time("date_birth").
			SchemaType(map[string]string{
				dialect.Postgres: "date",
			}),
	}
}
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("emails", Email.Type),
		edge.To("organization",Tenant.Type).Unique(),
	}
}
