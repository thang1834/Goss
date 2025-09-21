package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Cart struct {
	ent.Schema
}

func (Cart) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Cart) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("carts").Unique(),
		edge.To("items", CartItem.Type),
	}
}
