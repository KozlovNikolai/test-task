-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "products"(
    "id"  serial NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "provider_id" INTEGER NOT NULL,
    "price" FLOAT(53) NOT NULL,
    "stock" INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS "order_states"(
    "id"  serial NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS "users"(
    "id"  serial NOT NULL PRIMARY KEY,
    "login" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "token" TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS "providers"(
    "id"  serial NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "origin" TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS "items"(
    "id"  serial NOT NULL PRIMARY KEY,
    "product_id" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL,
    "total_price" FLOAT(53) NOT NULL,
    "order_id" INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS "orders"(
    "id"  serial NOT NULL PRIMARY KEY,
    "user_id" INTEGER NOT NULL,
    "state_id" INTEGER NOT NULL,
    "total_amount" FLOAT(53) NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "orders" ADD CONSTRAINT "orders_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "products" ADD CONSTRAINT "products_provider_id_foreign" FOREIGN KEY("provider_id") REFERENCES "providers"("id");
ALTER TABLE
    "items" ADD CONSTRAINT "items_product_id_foreign" FOREIGN KEY("product_id") REFERENCES "products"("id");
ALTER TABLE
    "orders" ADD CONSTRAINT "orders_state_id_foreign" FOREIGN KEY("state_id") REFERENCES "order_states"("id");
ALTER TABLE
    "items" ADD CONSTRAINT "items_order_id_foreign" FOREIGN KEY("order_id") REFERENCES "orders"("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table "items", "orders", "providers", "order_states", "users", "products";
-- +goose StatementEnd