-- +goose Up
-- +goose StatementBegin
INSERT INTO provider (id, name, origin) VALUES
(1, 'Star Inc', 'Kazakhstan'),
(2, 'Tesla', 'China'),
(3, 'Bosh', 'China'),
(4, 'West Wind Gmbh', 'Kazakhstan');

SELECT setval(pg_get_serial_sequence('provider', 'id'), max(id)) FROM provider;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM provider;
-- +goose StatementEnd