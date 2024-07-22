-- name: AddUser :exec
INSERT INTO users (
  email
) VALUES (
  $1
);

-- name: DeleteUser :exec
DELETE FROM users
WHERE email = $1;
