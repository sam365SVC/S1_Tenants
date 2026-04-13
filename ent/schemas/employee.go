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
		field.Enum("department").
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
		edge.From("branches", Branch.Type).Ref("employees").Unique(),
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

				departmentVal, existD := m.Field("department")
				posVal, existP := m.Field("position")

				// Si no se está tocando ninguno de los dos en un Update, saltar
				if !existD && !existP && !m.Op().Is(ent.OpCreate) {
					return next.Mutate(ctx, m)
				}

				var department, position string

				// Lógica para obtener los valores actuales (nuevos o viejos)
				if existD {
					department = fmt.Sprint(departmentVal)
				} else if !m.Op().Is(ent.OpCreate) {
					if oldDept, err := m.OldField(ctx, "department"); err == nil {
						department = fmt.Sprint(oldDept)
					}
				}

				if existP {
					position = posVal.(string)
				} else if !m.Op().Is(ent.OpCreate) {
					if oldPos, err := m.OldField(ctx, "position"); err == nil {
						position = fmt.Sprint(oldPos)
					}
				}

				// Validación
				validPositions, exists := allowedPositions[department]
				if !exists {
					return nil, fmt.Errorf("el área %s no es válida", department)
				}

				isAllowed := false
				for _, v := range validPositions {
					if v == position {
						isAllowed = true
						break
					}
				}

				if !isAllowed {
					return nil, fmt.Errorf("el cargo '%s' no está permitido para el área '%s'", position, department)
				}

				return next.Mutate(ctx, m)
			})
		},
	}
}
