-- name: AddEmail :exec
INSERT INTO sended_emails (
  email,
  updated_at
) VALUES (
  $1,
  $2
);

-- name: UpdateEmail :exec
UPDATE sended_emails
SET 
  updated_at = $2
WHERE
  email = $1;

-- name: DeleteEmail :exec
DELETE FROM sended_emails
WHERE email = $1;

-- name: GetAll :many
SELECT * FROM sended_emails;

