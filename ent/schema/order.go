package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Order struct {
	ent.Schema
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("status").Default("pending"),
		field.Float("total_price").Default(0),
		field.String("payment_method").Optional(),
		field.Text("shipping_address").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("orders").Unique(),
		edge.To("items", OrderItem.Type),
		edge.To("payments", Payment.Type),
	}
}
