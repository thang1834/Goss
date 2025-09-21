package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Category struct {
	ent.Schema
}

func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("name").MaxLen(100),
		field.String("slug").MaxLen(150).Unique(),
		field.Int("parent_id").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("products", Product.Type),
		edge.To("discounts", DiscountCategory.Type),
	}
}
