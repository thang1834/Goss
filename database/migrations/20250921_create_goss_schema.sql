-- +goose Up
-- +goose StatementBegin
-- Create table "users"
create table IF not exists  "users" (
    "id" BIGSERIAL PRIMARY KEY,
    "first_name" VARCHAR(255),
    "middle_name" VARCHAR(255),
    "last_name" VARCHAR(255),
    "email" VARCHAR(150) UNIQUE,
    "password_hash" VARCHAR(255),
    "phone" VARCHAR(20),
    "status" VARCHAR(20) DEFAULT 'active',
    "created_at" timestamp with time zone default current_timestamp,
    "updated_at" timestamp with time zone default current_timestamp,
    "verified_at" timestamp with time zone
);

-- Create table "roles"
create table IF not exists  "roles" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(50) UNIQUE,
    "description" TEXT,
    "is_active" BOOLEAN DEFAULT true,
    "created_at" timestamp with time zone default current_timestamp,
    "updated_at" timestamp with time zone default current_timestamp
);

-- Create table "permissions"
create table IF not exists  "permissions" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(100) UNIQUE,
    "description" TEXT,
    "resource" VARCHAR(50),
    "action" VARCHAR(20),
    "created_at" timestamp with time zone default current_timestamp,
    "updated_at" timestamp with time zone default current_timestamp
);

-- Create table "role_permissions"
create table IF not exists  "role_permissions" (
    "id" SERIAL PRIMARY KEY,
    "role_id" INTEGER,
    "permission_id" INTEGER,
    "created_at" timestamp with time zone default current_timestamp,
    CONSTRAINT "role_permissions_role_id_permission_id_key" UNIQUE ("role_id", "permission_id"),
    FOREIGN KEY ("role_id") REFERENCES "roles" ("id"),
    FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id")
);

-- Create table "user_roles"
create table IF not exists  "user_roles" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER,
    "role_id" INTEGER,
    "assigned_by" INTEGER,
    "assigned_at" timestamp with time zone default current_timestamp,
    "expires_at" TIMESTAMP,
    "is_active" BOOLEAN DEFAULT true,
    CONSTRAINT "user_roles_user_id_role_id_key" UNIQUE ("user_id", "role_id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
    FOREIGN KEY ("role_id") REFERENCES "roles" ("id"),
    FOREIGN KEY ("assigned_by") REFERENCES "users" ("id")
);

-- Create table "user_permissions"
create table IF not exists  "user_permissions" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER,
    "permission_id" INTEGER,
    "granted_by" INTEGER,
    "granted_at" timestamp with time zone default current_timestamp,
    "expires_at" TIMESTAMP,
    "is_active" BOOLEAN DEFAULT true,
    CONSTRAINT "user_permissions_user_id_permission_id_key" UNIQUE ("user_id", "permission_id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
    FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id"),
    FOREIGN KEY ("granted_by") REFERENCES "users" ("id")
);

-- Create table "categories"
create table IF not exists  "categories" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR(100),
    "slug" VARCHAR(150) UNIQUE,
    "parent_id" BIGINT,
    "created_at" timestamp with time zone default current_timestamp,
    "updated_at" timestamp with time zone default current_timestamp,
    FOREIGN KEY ("parent_id") REFERENCES "categories" ("id")
);

-- Create table "products"
create table IF not exists  "products" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR(200),
    "slug" VARCHAR(200) UNIQUE,
    "description" TEXT,
    "price" NUMERIC(12,2),
    "stock_quantity" INTEGER DEFAULT 0,
    "avg_rating" NUMERIC(3,2) DEFAULT 0,
    "review_count" INTEGER DEFAULT 0,
    "category_id" BIGINT,
    "created_at" timestamp with time zone default current_timestamp,
    "updated_at" timestamp with time zone default current_timestamp,
    FOREIGN KEY ("category_id") REFERENCES "categories" ("id")
);

-- Create table "product_images"
create table IF not exists  "product_images" (
    "id" BIGSERIAL PRIMARY KEY,
    "product_id" BIGINT,
    "image_url" TEXT,
    "is_primary" BOOLEAN DEFAULT false,
    FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

-- Create table "carts"
create table IF not exists  "carts" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT,
    "created_at" timestamp with time zone default current_timestamp,
    "updated_at" timestamp with time zone default current_timestamp,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- Create table "cart_items"
create table IF not exists  "cart_items" (
    "cart_id" BIGINT,
    "product_id" BIGINT,
    "quantity" INTEGER,
    PRIMARY KEY ("cart_id", "product_id"),
    FOREIGN KEY ("cart_id") REFERENCES "carts" ("id"),
    FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

-- Create table "orders"
create table IF not exists  "orders" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT,
    "status" VARCHAR(20) DEFAULT 'pending',
    "total_price" NUMERIC(12,2),
    "payment_method" VARCHAR(50),
    "shipping_address" TEXT,
    "created_at" timestamp with time zone default current_timestamp,
    "updated_at" timestamp with time zone default current_timestamp,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- Create table "order_items"
create table IF not exists  "order_items" (
    "id" BIGSERIAL PRIMARY KEY,
    "order_id" BIGINT,
    "product_id" BIGINT,
    "quantity" INTEGER,
    "unit_price" NUMERIC(12,2),
    FOREIGN KEY ("order_id") REFERENCES "orders" ("id"),
    FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

-- Create table "payments"
create table IF not exists  "payments" (
    "id" BIGSERIAL PRIMARY KEY,
    "order_id" BIGINT,
    "amount" NUMERIC(12,2),
    "method" VARCHAR(50),
    "status" VARCHAR(20) DEFAULT 'pending',
    "transaction_id" VARCHAR(100) UNIQUE,
    "created_at" timestamp with time zone default current_timestamp,
    FOREIGN KEY ("order_id") REFERENCES "orders" ("id")
);

-- Create table "reviews"
create table IF not exists  "reviews" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT,
    "product_id" BIGINT,
    "rating" INTEGER,
    "comment" TEXT,
    "created_at" timestamp with time zone default current_timestamp,
    "updated_at" timestamp with time zone default current_timestamp,
    CONSTRAINT "reviews_user_id_product_id_key" UNIQUE ("user_id", "product_id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
    FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

-- Create table "wishlists"
create table IF not exists  "wishlists" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT,
    "created_at" timestamp with time zone default current_timestamp,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);

-- Create table "wishlist_items"
create table IF not exists  "wishlist_items" (
    "wishlist_id" BIGINT,
    "product_id" BIGINT,
    "added_at" timestamp with time zone default current_timestamp,
    PRIMARY KEY ("wishlist_id", "product_id"),
    FOREIGN KEY ("wishlist_id") REFERENCES "wishlists" ("id"),
    FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

-- Create table "discounts"
create table IF not exists  "discounts" (
    "id" BIGSERIAL PRIMARY KEY,
    "code" VARCHAR(50) UNIQUE,
    "description" TEXT,
    "discount_type" VARCHAR(20),
    "discount_value" NUMERIC(12,2),
    "start_date" TIMESTAMP,
    "end_date" TIMESTAMP,
    "usage_limit" INTEGER,
    "usage_count" INTEGER DEFAULT 0,
    "min_order_value" NUMERIC(12,2),
    "created_at" timestamp with time zone default current_timestamp,
    "updated_at" timestamp with time zone default current_timestamp
);

-- Create table "discount_products"
create table IF not exists  "discount_products" (
    "discount_id" BIGINT,
    "product_id" BIGINT,
    PRIMARY KEY ("discount_id", "product_id"),
    FOREIGN KEY ("discount_id") REFERENCES "discounts" ("id"),
    FOREIGN KEY ("product_id") REFERENCES "products" ("id")
);

-- Create table "discount_categories"
create table IF not exists  "discount_categories" (
    "discount_id" BIGINT,
    "category_id" BIGINT,
    PRIMARY KEY ("discount_id", "category_id"),
    FOREIGN KEY ("discount_id") REFERENCES "discounts" ("id"),
    FOREIGN KEY ("category_id") REFERENCES "categories" ("id")
);

-- Create table "user_vouchers"
create table IF not exists  "user_vouchers" (
    "id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT,
    "discount_id" BIGINT,
    "is_used" BOOLEAN DEFAULT false,
    "used_at" TIMESTAMP,
    CONSTRAINT "user_vouchers_user_id_discount_id_key" UNIQUE ("user_id", "discount_id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
    FOREIGN KEY ("discount_id") REFERENCES "discounts" ("id")
);

CREATE TABLE IF NOT EXISTS sessions
(
    token   TEXT PRIMARY KEY,
    user_id BIGINT      CONSTRAINT session_user_fk REFERENCES users ON DELETE CASCADE ,
    data    BYTEA       NOT NULL,
    expiry  TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user_vouchers";
DROP TABLE "discount_categories";
DROP TABLE "discount_products";
DROP TABLE "discounts";
DROP TABLE "wishlist_items";
DROP TABLE "wishlists";
DROP TABLE "reviews";
DROP TABLE "payments";
DROP TABLE "order_items";
DROP TABLE "orders";
DROP TABLE "cart_items";
DROP TABLE "carts";
DROP TABLE "product_images";
DROP TABLE "products";
DROP TABLE "categories";
DROP TABLE "user_permissions";
DROP TABLE "user_roles";
DROP TABLE "role_permissions";
DROP TABLE "permissions";
DROP TABLE "roles";
DROP TABLE "users";
DROP TABLE sessions;
-- +goose StatementEnd