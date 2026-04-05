package schemas

import (
	"errors"
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}
var nameRegex=regexp.MustCompile(`^[A-ZÁÉÍÓÚÑ][a-záéíóúñ]+(?:\s[A-ZÁÉÍÓÚÑ][a-záéíóúñ]+)*$`)
var emailRegex = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
func (User)Fields()[]ent.Field  {
	return []ent.Field{
		field.String("name").MaxLen(50).
			Validate(func(s string) error {
				if !nameRegex.MatchString(s) {
					return errors.New("Each word in the name must start with a capital letter, and numbers are not allowed")
				}
				return nil
			}).
			NotEmpty(),
		field.String("last_name").MaxLen(50).
			Validate(func(s string) error {
				if !nameRegex.MatchString(s) {
					return errors.New("Each word in the last name must start with a capital letter, and numbers are not allowed")
				}
				return nil
			}).
			NotEmpty(),
		field.Int("ci").Range(11111,999999999).Unique(),
		field.Enum("rol").Values("DEVELOPER","USER").Default("USER"),
		field.String("phone").MaxLen(14).NotEmpty(),
		field.Time("date_birth").
			SchemaType(map[string]string{
				dialect.Postgres:"date",
			}),
		field.String("email").MaxLen(100).Unique().Unique(),
		field.String("password_hash").MaxLen(150),
	}
}
func (User) Edges()[]ent.Edge{
	return []ent.Edge{
		edge.To("employees",Employee.Type),
	}
}