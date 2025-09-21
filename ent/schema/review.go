package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Review struct {
	ent.Schema
}

func (Review) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Int("rating"),
		field.Text("comment").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Review) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("reviews").Unique(),
		edge.From("product", Product.Type).Ref("reviews").Unique(),
	}
}

//func (Review) Indexes() []ent.Index {
//	return []ent.Index{
//		index.Fields("user_id", "product_id").Unique(),
//	}
//}
