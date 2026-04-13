package schemas

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)
func defaultExpireAt() time.Time {
    return time.Now().Add(2 * time.Hour).Add(15 * time.Minute)
}
type Invitation struct {
	ent.Schema
}

func (Invitation) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").NotEmpty().MaxLen(100).Unique().Validate(ValidateEmail),
		field.String("token").NotEmpty().MaxLen(64).Unique(),
		field.Enum("account").Values("DEVELOPER","USER","ADMIN"),
		field.Time("expire_at").
			Default(defaultExpireAt()),
	}
}
func (Invitation) Edges()[]ent.Edge {
	return []ent.Edge{
		edge.To("invitation_employee",InvitationEmployee.Type).Unique(),
	}
}
