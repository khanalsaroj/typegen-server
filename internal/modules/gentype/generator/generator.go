package gen

import (
	"database/sql"

	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type Generator interface {
	Generate(rows *sql.Rows, req domain.TypeRequest, tbN string, dbType string) (string, error)
}
