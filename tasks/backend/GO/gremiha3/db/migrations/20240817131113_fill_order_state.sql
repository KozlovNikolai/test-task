-- +goose Up
-- +goose StatementBegin
INSERT INTO order_state (id, name) VALUES
(1, 'created'),
(2, 'in process'),
(3, 'pending'),
(4, 'paid');

SELECT setval(pg_get_serial_sequence('order_state', 'id'), max(id)) FROM order_state;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM order_state;
-- +goose StatementEnd