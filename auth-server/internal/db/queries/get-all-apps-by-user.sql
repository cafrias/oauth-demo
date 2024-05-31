-- name: GetAllAppsByUser :many
SELECT * FROM app WHERE UserID = ? ORDER BY CREATED ASC;

