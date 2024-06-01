-- name: ResetAppSecret :exec
UPDATE app SET Hash=? WHERE ClientID=? AND UserID=?;
