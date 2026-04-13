package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Plan struct {
	ent.Schema
}

func (Plan) Fields() []ent.Field {
	return []ent.Field{
		field.String("subscription").MaxLen(20),
		field.Float("price").
			SchemaType(map[string]string{
				dialect.Postgres: "decimal(7,2)",
			}).
			Default(0.00),
		field.Int32("max_employees").Default(5),
		field.Int32("max_branches").Default(2),
		field.Int32("max_bosses").Default(1),
	}
}
func (Plan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tenants", Tenant.Type),
	}
}
