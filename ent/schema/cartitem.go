package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type CartItem struct {
	ent.Schema
}

func (CartItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("quantity"),
	}
}

func (CartItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cart", Cart.Type).Ref("items").Unique(),
		edge.From("product", Product.Type).Ref("cart_items").Unique(),
	}
}

//func (CartItem) Indexes() []ent.Index {
//	return []ent.Index{
//		index.Fields("cart_id", "product_id").Unique(),
//	}
//}
