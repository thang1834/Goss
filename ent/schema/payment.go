package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Payment struct {
	ent.Schema
}

func (Payment) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Float("amount"),
		field.String("method").Optional(),
		field.String("status").Default("pending"),
		field.String("transaction_id").Optional(),
		field.Time("created_at").Default(time.Now),
	}
}

func (Payment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).Ref("payments").Unique(),
	}
}
