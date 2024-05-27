-- +goose Up
-- +goose StatementBegin
CREATE TABLE user (
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	Email VARCHAR(255) NOT NULL,
	Hash VARCHAR(255) NOT NULL,
	Created DATETIME DEFAULT CURRENT_TIMESTAMP,
	Updated DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX user_email ON user (Email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX user_email;
DROP TABLE user;
-- +goose StatementEnd
