package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Email struct {
	ent.Schema
}
func (Email)Fields()[]ent.Field{
	return []ent.Field{
		field.String("email").NotEmpty().MaxLen(150).Validate(ValidateEmail).Unique(),
		field.String("pasword_hash").NotEmpty().MaxLen(250),
		field.Enum("account").Values("DEVELOPER","USER","ADMIN"),
		field.Time("create_at").Default(time.Now).Immutable(),
		field.Time("update_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Email)Edges()[]ent.Edge{
	return []ent.Edge{
		edge.From("user",User.Type).Ref("emails").Unique(),
		edge.To("employees",Employee.Type),
	}
}


