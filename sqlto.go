// sqlto package provides methods to take database/sql rows and convert them to
// different export formats.
package sqlto

import "database/sql"

type SQLto struct {
	Rows *sql.Rows
}

func New(rows *sql.Rows) *SQLto {

	sqlto := SQLto{}
	sqlto.Rows = rows

	return &sqlto
}
