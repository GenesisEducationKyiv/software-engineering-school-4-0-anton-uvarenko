// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: sended_emails.sql

package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addEmail = `-- name: AddEmail :exec
INSERT INTO sended_emails (
  email,
  updated_at
) VALUES (
  $1,
  $2
)
`

type AddEmailParams struct {
	Email     pgtype.Text
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) AddEmail(ctx context.Context, arg AddEmailParams) error {
	_, err := q.db.Exec(ctx, addEmail, arg.Email, arg.UpdatedAt)
	return err
}

const deleteEmail = `-- name: DeleteEmail :exec
DELETE FROM sended_emails
WHERE email = $1
`

func (q *Queries) DeleteEmail(ctx context.Context, email pgtype.Text) error {
	_, err := q.db.Exec(ctx, deleteEmail, email)
	return err
}

const getAll = `-- name: GetAll :many
SELECT email, updated_at FROM sended_emails
`

func (q *Queries) GetAll(ctx context.Context) ([]SendedEmail, error) {
	rows, err := q.db.Query(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SendedEmail
	for rows.Next() {
		var i SendedEmail
		if err := rows.Scan(&i.Email, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateEmail = `-- name: UpdateEmail :exec
UPDATE sended_emails
SET 
  updated_at = $2
WHERE
  email = $1
`

type UpdateEmailParams struct {
	Email     pgtype.Text
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) UpdateEmail(ctx context.Context, arg UpdateEmailParams) error {
	_, err := q.db.Exec(ctx, updateEmail, arg.Email, arg.UpdatedAt)
	return err
}
