package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type WishlistItem struct {
	ent.Schema
}

func (WishlistItem) Fields() []ent.Field {
	return []ent.Field{
		field.Time("added_at").Default(time.Now),
	}
}

func (WishlistItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("wishlist", Wishlist.Type).Ref("items").Unique(),
		edge.From("product", Product.Type).Ref("wishlist_items").Unique(),
	}
}

//func (WishlistItem) Indexes() []ent.Index {
//	return []ent.Index{
//		index.Fields("wishlist_id", "product_id").Unique(),
//	}
//}
