-- +goose Up
-- +goose StatementBegin
INSERT INTO order_states (id, name) VALUES
(1, 'created'),
(2, 'in process'),
(3, 'pending'),
(4, 'paid');

SELECT setval(pg_get_serial_sequence('order_states', 'id'), max(id)) FROM order_states;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM order_states;
-- +goose StatementEnd