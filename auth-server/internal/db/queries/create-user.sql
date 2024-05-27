-- name: CreateUser :one
INSERT INTO user (
	Email,
	Hash
) VALUES (
	?,
	?
) RETURNING id;
