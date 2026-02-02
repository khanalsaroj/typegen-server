package helper

import (
	"database/sql"
	"errors"

	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/infrastructure/db/connector"
)

func OpenDatabase(connDetails *domain.DatabaseConnection) (*sql.DB, connector.DBConnector, domain.DatabaseConnectionInfo, error) {
	if connDetails == nil {
		return nil, nil, domain.DatabaseConnectionInfo{}, errors.New("connection details are nil")
	}

	connInfo := convertDatabaseConnection(connDetails)

	conn, err := connector.New(connInfo)
	if err != nil {
		return nil, nil, connInfo, err
	}

	db, err := conn.Open(connInfo)
	if err != nil {
		return nil, nil, connInfo, err
	}

	reader, err := connector.New(connInfo)
	if err != nil {
		if db.Close() != nil {
			return nil, nil, domain.DatabaseConnectionInfo{}, err
		}
		return nil, nil, connInfo, err
	}
	return db, reader, connInfo, nil
}

func convertDatabaseConnection(conn *domain.DatabaseConnection) domain.DatabaseConnectionInfo {
	return domain.DatabaseConnectionInfo{
		DbType:       conn.DbType,
		Host:         conn.Host,
		Port:         conn.Port,
		Username:     conn.Username,
		Password:     conn.Password,
		SchemaName:   conn.SchemaName,
		DatabaseName: conn.DatabaseName,
	}
}
