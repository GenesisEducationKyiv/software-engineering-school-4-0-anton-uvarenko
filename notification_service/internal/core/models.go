// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package repo

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type SendedEmail struct {
	Email     pgtype.Text
	UpdatedAt pgtype.Timestamp
}