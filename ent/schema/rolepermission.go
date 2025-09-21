package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type RolePermission struct {
	ent.Schema
}

func (RolePermission) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Uint64("role_id"),
		field.Uint64("permission_id"),
		field.Time("created_at").Default(time.Now),
	}
}

func (RolePermission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", Role.Type).
			Ref("permissions").
			Field("role_id").
			Required().
			Unique(),
		edge.From("permission", Permission.Type).
			Ref("roles").
			Field("permission_id").
			Required().
			Unique(),
	}
}

func (RolePermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id", "permission_id").Unique(),
	}
}
