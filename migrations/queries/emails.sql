-- name: AddEmail :exec
INSERT INTO emails (
  email
) VALUES (
  $1
);

-- name: GetAll :many
SELECT * FROM emails;
