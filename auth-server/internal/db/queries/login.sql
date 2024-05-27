-- name: Login :one
SELECT * FROM user WHERE Email = ? AND Hash = ?;

