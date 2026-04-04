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

func (Employee)Fields()[]ent.Field{
	return []ent.Field{
		field.Enum("area").
			Values("office","logistics","plant","commercial"),
		field.String("position").MaxLen(30),
		field.Bool("active").Default(true),
		field.Time("join_at").Default(time.Now),
		field.Time("left_at").Optional().Nillable(),
	}
}
func (Employee) Edges()[]ent.Edge{
	return []ent.Edge{
		edge.From("user",User.Type).Ref("employees").Unique(),
		edge.From("tenant",Tenant.Type).Ref("employees").Unique(),
		edge.From("branch",Branch.Type).Ref("employees").Unique(),
	}
}

func (Employee) Hooks()[]ent.Hook{
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if (!m.Op().Is(ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne)) {
					return next.Mutate(ctx,m) 
				}
				areaVal,existA:=m.Field("area")
				positionVal,existP:=m.Field("position")

				if (!existA&&!existP) {
					return next.Mutate(ctx,m)
				}
				var area,position string

				oldValueGettes,ok:=m.(interface{
					OldArea(context.Context)(string,error)
					OldPosition(context.Context)(string,error)
				})
				if existA {
					area=areaVal.(string)
				}else if ok {
					area,_=oldValueGettes.OldArea(ctx)
				}

				if existP {
					position=positionVal.(string)
				}else if ok {
					position,_=oldValueGettes.OldPosition(ctx)
				}

				allowedPositions:=map[string][]string{
					"office":	{"boss","manager","admin","accountant"},
					"logistics":{"technician","driver","dispatcher"},
					"plant":	{"supervisor","operator","sorter"},
					"commercial":{"sales_rep","manager_commercial"},
				}
				validationPosition,exists:=allowedPositions[area]
				if !exists {
					return nil,fmt.Errorf("el area %s no es valida",area)
				}
				isAllowed:=false

				for _, v := range validationPosition {
					if v==position {
						isAllowed=true
						break
					}
				}
				if !isAllowed {
					return nil,fmt.Errorf("el cargo %s no ésta permitido para el área %s",position,area)
				}
				return next.Mutate(ctx,m)
			})
		},
	}
}