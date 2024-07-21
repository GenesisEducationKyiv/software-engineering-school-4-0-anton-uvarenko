-- name: AddEmail :exec
INSERT INTO sened_emails (
  email
  updated_at
) VALUES (
  $1,
  $2
);

-- name: GetAll :many
SELECT * FROM emails;

