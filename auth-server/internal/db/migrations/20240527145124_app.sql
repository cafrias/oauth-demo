-- +goose Up
-- +goose StatementBegin
CREATE TABLE app (
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	Name VARCHAR(255) NOT NULL,
	UserID INTEGER NOT NULL,
	Type VARCHAR(255) NOT NULL,
	RedirectURI VARCHAR(255) NOT NULL,
	ClientID CHAR(20) NOT NULL,
	Hash VARCHAR(255) NOT NULL,
	Created DATETIME DEFAULT CURRENT_TIMESTAMP,
	Updated DATETIME DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (UserID) REFERENCES user(ID)
);

CREATE UNIQUE INDEX app_client_id ON app (ClientID);
CREATE INDEX app_user_id ON app (UserID);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX app_client_id;
DROP INDEX app_user_id;
DROP TABLE app;
-- +goose StatementEnd
