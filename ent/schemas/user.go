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

func (User)Fields()[]ent.Field  {
	return []ent.Field{
		field.String("name").MaxLen(50).NotEmpty(),
		field.String("last_name").MaxLen(50).NotEmpty(),
		field.Int("ci").Range(11111,99999999),
		field.Time("date_birth").
			SchemaType(map[string]string{
				dialect.Postgres:"date",
			}),
		field.String("email").MaxLen(100).Unique(),
		field.String("password_hash").MaxLen(150),
	}
}
func (User) Edges()[]ent.Edge{
	return []ent.Edge{
		edge.To("employees",Employee.Type),
	}
}