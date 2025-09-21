package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Permission struct {
	ent.Schema
}

func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("name").MaxLen(255).Unique(),
		field.Text("description").Optional(),
		field.String("resource").MaxLen(255),
		field.String("action").MaxLen(255),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", RolePermission.Type),
		edge.To("user_permissions", UserPermission.Type),
	}
}
