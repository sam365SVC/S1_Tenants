package schemas

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Employee struct {
	ent.Schema
}

func (Employee) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("area").
			Values("office", "logistics", "plant", "commercial"),
		field.String("position").MaxLen(30),
		field.Bool("active").Default(true),
		field.Time("join_at").Default(time.Now),
		field.Time("left_at").Optional().Nillable(),
	}
}
func (Employee) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("emails", Email.Type).Ref("employees").Unique(),
		edge.From("tenant", Tenant.Type).Ref("employees").Unique(),
		edge.From("branch", Branch.Type).Ref("employees").Unique(),
	}
}

var allowedPositions = map[string][]string{
	"office":     {"boss", "manager", "admin", "accountant"},
	"logistics":  {"technician", "driver", "dispatcher"},
	"plant":      {"supervisor", "operator", "sorter"},
	"commercial": {"sales_rep", "manager_commercial"},
}

func (Employee) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				// Solo validamos en creación y edición
				if !m.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne) {
					return next.Mutate(ctx, m)
				}

				areaVal, existA := m.Field("area")
				posVal, existP := m.Field("position")

				// Si no se está tocando ninguno de los dos en un Update, saltar
				if !existA && !existP && !m.Op().Is(ent.OpCreate) {
					return next.Mutate(ctx, m)
				}

				var area, position string

				// Lógica para obtener los valores actuales (nuevos o viejos)
				if existA {
					area = areaVal.(string)
				} else {
					// Si es Update y no viene 'area', buscamos el valor que ya tenía
					if getter, ok := m.(interface {
						OldArea(context.Context) (string, error)
					}); ok {
						area, _ = getter.OldArea(ctx)
					}
				}

				if existP {
					position = posVal.(string)
				} else {
					// Si es Update y no viene 'position', buscamos el valor que ya tenía
					if getter, ok := m.(interface {
						OldPosition(context.Context) (string, error)
					}); ok {
						position, _ = getter.OldPosition(ctx)
					}
				}

				// Validación
				validPositions, exists := allowedPositions[area]
				if !exists {
					return nil, fmt.Errorf("el área %s no es válida", area)
				}

				isAllowed := false
				for _, v := range validPositions {
					if v == position {
						isAllowed = true
						break
					}
				}

				if !isAllowed {
					return nil, fmt.Errorf("el cargo '%s' no está permitido para el área '%s'", position, area)
				}

				return next.Mutate(ctx, m)
			})
		},
	}
}
