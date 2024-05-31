-- name: FindUserByEmail :one
SELECT * FROM user WHERE Email = ?;

