package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type UserPermission struct {
	ent.Schema
}

func (UserPermission) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Uint64("user_id"),
		field.Uint64("permission_id"),
		field.Uint64("granted_by"),
		field.Time("granted_at").Default(time.Now),
		field.Time("expires_at").Optional(),
		field.Bool("is_active").Default(true),
	}
}

func (UserPermission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("permissions").
			Field("user_id").
			Required().
			Unique(),
		edge.From("permission", Permission.Type).
			Ref("user_permissions").
			Field("permission_id").
			Required().
			Unique(),
	}
}

func (UserPermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "permission_id").Unique(),
	}
}
