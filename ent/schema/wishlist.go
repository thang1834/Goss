package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Wishlist struct {
	ent.Schema
}

func (Wishlist) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Time("created_at").Default(time.Now),
	}
}

func (Wishlist) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("wishlists").Unique(),
		edge.To("items", WishlistItem.Type),
	}
}
