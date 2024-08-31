-- +goose Up
-- +goose StatementBegin
INSERT INTO products (id, name, provider_id, price, stock) VALUES
(1, 'телефон', 1, 12300.00, 113),
(2, 'пылесос', 2, 8750.35, 76),
(3, 'кофеварка', 3, 4365.65, 32),
(4, 'утюг', 4, 3276.15, 75),
(5, 'яхта', 1, 4567983.23, 1),
(6, 'автомобиль', 2, 2134678.23, 3),
(7, 'ноутбук', 3, 87456.34, 25),
(8, 'принтер', 4, 19456.32, 15),
(9, 'стол', 1, 3546.87, 25),
(10, 'стул', 2, 1254.65, 100),
(11, 'веник', 3, 234.87, 100),
(12, 'телевизор', 4, 34876.31, 42),
(13, 'лампа', 1, 954.93, 30),
(14, 'шкаф', 2, 12876.89, 5),
(15, 'шапка', 3, 123.76, 250);

SELECT setval(pg_get_serial_sequence('products', 'id'), max(id)) FROM products;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM products;
-- +goose StatementEnd