package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type OrderItem struct {
	ent.Schema
}

func (OrderItem) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Int("quantity"),
		field.Float("unit_price"),
	}
}

func (OrderItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("items").Unique(),
		edge.From("product", Product.Type).Ref("order_items").Unique(),
	}
}
