-- name: CreateApp :one
INSERT INTO app (
	Name,
	UserID,
	Type,
	RedirectURI,
	ClientID,
	Hash
) VALUES(
	?,
	?,
	?,
	?,
	?,
	?
) RETURNING id;
