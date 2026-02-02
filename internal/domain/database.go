package domain

import (
	"time"
)

type DatabaseConnection struct {
	ConnectionID uint64    `gorm:"column:connection_id;primaryKey;autoIncrement" json:"connectionId"`
	Name         string    `gorm:"column:name;size:100;not null" json:"name"`
	DbType       string    `gorm:"column:db_type;size:20;not null" json:"dbType"`
	Host         string    `gorm:"column:host;size:255;not null" json:"host"`
	Port         int       `gorm:"column:port;not null" json:"port"`
	DatabaseName string    `gorm:"column:database_name;size:100;not null" json:"databaseName"`
	SchemaName   string    `gorm:"column:schema_name;size:100" json:"schemaName,omitempty"`
	Username     string    `gorm:"column:username;size:100;not null" json:"username"`
	Password     string    `gorm:"column:password;type:text;not null" json:"password"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (u *DatabaseConnection) TableName() string {
	return "database_connections"
}
