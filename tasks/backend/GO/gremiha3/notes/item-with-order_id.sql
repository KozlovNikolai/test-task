-- +goose Up
-- +goose StatementBegin
CREATE TABLE "product"(
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "provider_id" INTEGER NOT NULL,
    "price" FLOAT(53) NOT NULL,
    "stock" INTEGER NOT NULL
);
ALTER TABLE
    "product" ADD PRIMARY KEY("id");
CREATE TABLE "order_state"(
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL
);
ALTER TABLE
    "order_state" ADD PRIMARY KEY("id");
CREATE TABLE "user"(
    "id" INTEGER NOT NULL,
    "login" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "token" TEXT NOT NULL
);
ALTER TABLE
    "user" ADD PRIMARY KEY("id");
CREATE TABLE "provider"(
    "id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "origin" TEXT NOT NULL
);
ALTER TABLE
    "provider" ADD PRIMARY KEY("id");
CREATE TABLE "item"(
    "id" INTEGER NOT NULL,
    "product_id" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL,
    "total_price" FLOAT(53) NOT NULL,
    "order_id" INTEGER NOT NULL
);
ALTER TABLE
    "item" ADD PRIMARY KEY("id");
CREATE TABLE "order"(
    "id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "state_id" INTEGER NOT NULL,
    "total_amount" FLOAT(53) NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "order" ADD PRIMARY KEY("id");
ALTER TABLE
    "order" ADD CONSTRAINT "order_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "user"("id");
ALTER TABLE
    "product" ADD CONSTRAINT "product_provider_id_foreign" FOREIGN KEY("provider_id") REFERENCES "provider"("id");
ALTER TABLE
    "item" ADD CONSTRAINT "item_product_id_foreign" FOREIGN KEY("product_id") REFERENCES "product"("id");
ALTER TABLE
    "order" ADD CONSTRAINT "order_state_id_foreign" FOREIGN KEY("state_id") REFERENCES "order_state"("id");
ALTER TABLE
    "item" ADD CONSTRAINT "item_order_id_foreign" FOREIGN KEY("order_id") REFERENCES "order"("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table "item", "order", "provider", "order_state", "user", "product";
-- +goose StatementEnd