-- name: AddEmail :exec
INSERT INTO sended_emails (
  email,
  updated_at
) VALUES (
  $1,
  $2
);

-- name: GetAll :many
SELECT * FROM sended_emails;

