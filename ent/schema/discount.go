package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Discount struct {
	ent.Schema
}

func (Discount) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("code").Unique(),
		field.Text("description").Optional(),
		field.String("discount_type"),
		field.Float("discount_value"),
		field.Time("start_date"),
		field.Time("end_date"),
		field.Int("usage_limit").Optional(),
		field.Int("usage_count").Default(0),
		field.Int("min_order_value").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Discount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("products", DiscountProduct.Type),
		edge.To("categories", DiscountCategory.Type),
		edge.To("user_vouchers", UserVoucher.Type),
	}
}
