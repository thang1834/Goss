package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ProductImage struct {
	ent.Schema
}

func (ProductImage) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("image_url"),
		field.Bool("is_primary").Default(false),
	}
}

func (ProductImage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("product", Product.Type).Ref("images").Unique(),
	}
}
