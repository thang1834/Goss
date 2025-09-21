package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

type DiscountProduct struct {
	ent.Schema
}

func (DiscountProduct) Fields() []ent.Field { return nil }

func (DiscountProduct) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("discount", Discount.Type).Ref("products").Unique(),
		edge.From("product", Product.Type).Ref("discounts").Unique(),
	}
}

//func (DiscountProduct) Indexes() []ent.Index {
//	return []ent.Index{
//		index.Fields("discount_id", "product_id").Unique(),
//	}
//}
