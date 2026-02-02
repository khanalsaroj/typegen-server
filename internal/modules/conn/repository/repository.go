package repository

import (
	"database/sql"
	"time"

	"github.com/khanalsaroj/typegen-server/internal/infrastructure/db/connector"

	"github.com/khanalsaroj/typegen-server/internal/domain"
)

type Repo interface {
	PingDB(req domain.DatabaseConnectionInfo) (
		pingMs int64,
		tablesFound int,
		sizeMb float64,
		err error,
		tableInfo []domain.TableInfo,
	)
}

type repo struct{}

func New() Repo {
	return &repo{}
}

func (r *repo) PingDB(req domain.DatabaseConnectionInfo) (
	int64,
	int,
	float64,
	error,
	[]domain.TableInfo,
) {
	conn, err := connector.New(req)
	if err != nil {
		return 0, 0, 0, err, nil
	}

	db, err := conn.Open(req)
	if err != nil {
		return 0, 0, 0, err, nil
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	start := time.Now()
	if err := db.Ping(); err != nil {
		return 0, 0, 0, err, nil
	}
	pingMs := time.Since(start).Milliseconds()

	stats, err := conn.Stats(db, req.DatabaseName, req.SchemaName)
	if err != nil {
		return 0, 0, 0, err, nil
	}

	return pingMs, stats.Tables, stats.SizeMB, nil, stats.TablesInfo
}
