package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the Author.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("first_name").Optional().MaxLen(255),
		field.String("middle_name").Optional().MaxLen(255),
		field.String("last_name").Optional().MaxLen(255),
		field.String("email").MaxLen(255).Unique(),
		field.String("password_hash").MaxLen(255),
		field.String("phone").Optional().MaxLen(20),
		field.String("status").MaxLen(50).Default("active"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("verified_at").Optional().Nillable().StructTag(`json:"-"`),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("roles", UserRole.Type),
		edge.To("permissions", UserPermission.Type),
		edge.To("carts", Cart.Type),
		edge.To("orders", Order.Type),
		edge.To("wishlists", Wishlist.Type),
		edge.To("reviews", Review.Type),
		edge.To("vouchers", UserVoucher.Type),
	}
}
