package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type UserRole struct {
	ent.Schema
}

func (UserRole) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Uint64("user_id"),
		field.Uint64("role_id"),
		field.Uint64("assigned_by"),
		field.Time("assigned_at").Default(time.Now),
		field.Time("expires_at").Optional(),
		field.Bool("is_active").Default(true),
	}
}

func (UserRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("roles").
			Field("user_id").
			Required().
			Unique(),
		edge.From("role", Role.Type).
			Ref("user_roles").
			Field("role_id").
			Required().
			Unique(),
	}
}

func (UserRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "role_id").Unique(),
	}
}
