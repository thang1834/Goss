package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Product struct {
	ent.Schema
}

func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("name").MaxLen(200),
		field.String("slug").MaxLen(200).Unique(),
		field.Text("description").Optional(),
		field.Float("price"),
		field.Int("stock_quantity").Default(0),
		field.Float("avg_rating").Default(0),
		field.Int("review_count").Default(0),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Product) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", Category.Type).Ref("products").Unique(),
		edge.To("images", ProductImage.Type),
		edge.To("reviews", Review.Type),
		edge.To("cart_items", CartItem.Type),
		edge.To("order_items", OrderItem.Type),
		edge.To("discounts", DiscountProduct.Type),
		edge.To("wishlist_items", WishlistItem.Type),
	}
}
