-- name: AddUser :exec
INSERT INTO users (
  email
) VALUES (
  $1
);
