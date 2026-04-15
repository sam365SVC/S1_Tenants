package schemas
import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Branch struct {
	ent.Schema
}
func (Branch) Fields() []ent.Field {
	return []ent.Field{
		// Datos Generales
		field.String("name").NotEmpty().MaxLen(100).Unique(),
		field.String("phone").Optional().MaxLen(14).Unique(), // Optional, por si es un acopio sin teléfono fijo

		// Dirección Estructurada
		field.String("street").
			NotEmpty().
			MaxLen(150).
			Comment("Avenida o calle principal"),
		field.String("zone").
			Optional().
			MaxLen(100).
			Comment("Barrio, UV o Zona"),
		field.String("reference").
			Optional().
			MaxLen(200).
			Comment("Ej: A dos cuadras del surtidor, portón verde"),
		
		field.String("city").
			NotEmpty().
			MaxLen(50).
			Default("Santa Cruz de la Sierra"), // Puedes poner un default útil
		
		// Departamentos completos de Bolivia
		field.Enum("state").
			Values("LP", "SC", "CB", "OR", "PT", "TJ", "CH", "BE", "PA").
			Default("LP"),

		// si es que quiere añadir mapas donde la gente pueda ver las dirrecciones mediante un mapa creado por el frontend
		// field.Float("latitude").Optional().Nillable(),
		// field.Float("longitude").Optional().Nillable(),

		// Metadatos
		field.Bool("is_active").Default(true),
		field.Time("create_at").Default(time.Now).Immutable(), // Immutable: no se puede editar esta fecha después de creada
	}
}

func (Branch) Edges()[]ent.Edge{
	return []ent.Edge{
		edge.From("tenant",Tenant.Type).Ref("branches").Unique(),
		edge.To("employees",Employee.Type),
		edge.To("invitation_employees",InvitationEmployee.Type),
	}
}
