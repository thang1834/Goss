package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

type DiscountCategory struct {
	ent.Schema
}

func (DiscountCategory) Fields() []ent.Field { return nil }

func (DiscountCategory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("discount", Discount.Type).Ref("categories").Unique(),
		edge.From("category", Category.Type).Ref("discounts").Unique(),
	}
}

//func (DiscountCategory) Indexes() []ent.Index {
//	return []ent.Index{
//		index.Fields("discount_id", "category_id").Unique(),
//	}
//}
