package generator

import (
	"database/sql"
	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type Mapper interface {
	Generate(rows *sql.Rows, req domain.MapperRequest, tbN string) (string, error)
}
