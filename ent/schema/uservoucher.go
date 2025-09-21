package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type UserVoucher struct {
	ent.Schema
}

func (UserVoucher) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Bool("is_used").Default(false),
		field.Time("used_at").Optional(),
	}
}

func (UserVoucher) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("vouchers").Unique(),
		edge.From("discount", Discount.Type).Ref("user_vouchers").Unique(),
	}
}

//func (UserVoucher) Indexes() []ent.Index {
//	return []ent.Index{
//		index.Fields("user_id", "discount_id").Unique(),
//	}
//}
