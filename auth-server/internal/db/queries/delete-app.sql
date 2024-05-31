-- name: DeleteApp :exec
DELETE FROM app WHERE ClientID = ? AND UserID = ?;

