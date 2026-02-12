package domain

import (
	"database/sql"
	"encoding/json"
	"time"
)

type HealthStatus string

type TableInfo struct {
	Name        string `json:"name"`
	ColumnCount int    `json:"columnCount"`
}

type ConnectResponse struct {
	ConnectionId int         `json:"connectionId"`
	Message      string      `json:"message"`
	Success      bool        `json:"success"`
	PingMs       int64       `json:"pingMs,omitempty"`
	TablesFound  int         `json:"tablesFound,omitempty"`
	SizeMb       float64     `json:"sizeMb,omitempty"`
	TableInfo    []TableInfo `json:"tables"`
}

type SqlData struct {
	Ordinal                int            `json:"ordinal"`
	ColumnName             string         `json:"columnName"`
	IsNullable             string         `json:"isNullable"`
	CharacterMaximumLength sql.NullInt16  `json:"characterMaximumLength"`
	DataType               string         `json:"dataType"`
	ColumnKey              string         `json:"columnKey"`
	ColumnComment          sql.NullString `json:"columnComment"`
}

type Stats struct {
	Tables     int
	SizeMB     float64
	TablesInfo []TableInfo
}

type DatabaseConnectionInfo struct {
	DbType       string `json:"dbType"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	SchemaName   string `json:"schemaName,omitempty"`
	DatabaseName string `json:"databaseName"`
}

type TypeRequest struct {
	ConnectionId   uint            `json:"connectionId"`
	Options        json.RawMessage `json:"options"`
	Prefix         string          `json:"prefix,omitempty"`
	Suffix         string          `json:"suffix,omitempty"`
	Style          string          `json:"style,omitempty"`
	TargetLanguage string          `json:"language"`
	TableNames     []string        `json:"tableNames"`
}

type MapperRequest struct {
	ConnectionId uint            `json:"connectionId"`
	Options      json.RawMessage `json:"options"`
	TargetType   string          `json:"targetType"`
	TableName    string          `json:"tableName"`
}

type JavaOptions struct {
	Getter             bool `json:"getter,omitempty"`
	Setter             bool `json:"setter,omitempty"`
	NoArgsConstructor  bool `json:"noArgsConstructor,omitempty"`
	AllArgsConstructor bool `json:"allArgsConstructor,omitempty"`
	Builder            bool `json:"builder,omitempty"`
	Data               bool `json:"data,omitempty"`
	SwaggerAnnotations bool `json:"swaggerAnnotations,omitempty"`
	Serializable       bool `json:"serializable,omitempty"`
	JacksonAnnotations bool `json:"jacksonAnnotations,omitempty"`
	ExtraSpacing       bool `json:"extraSpacing,omitempty"`
}

type RecordOptions struct {
	SwaggerAnnotations bool `json:"swaggerAnnotations,omitempty"`
	JacksonAnnotations bool `json:"jacksonAnnotations,omitempty"`
	Builder            bool `json:"builder,omitempty"`
	ExtraSpacing       bool `json:"extraSpacing,omitempty"`
}

type CSharpDtoOptions struct {
	ExtraSpacing        bool `json:"extraSpacing,omitempty"`
	CamelCaseProperties bool `json:"camelCaseProperties,omitempty"`
	Nullable            bool `json:"nullable,omitempty"`
	Getter              bool `json:"getter,omitempty"`
	Setter              bool `json:"setter,omitempty"`
	JsonPropertyName    bool `json:"jsonPropertyName,omitempty"`
}

type CSharpRecordOptions struct {
	CamelCaseProperties bool `json:"camelCaseProperties,omitempty"`
	ExtraSpacing        bool `json:"extraSpacing,omitempty"`
	Nullable            bool `json:"nullable,omitempty"`
	JsonPropertyName    bool `json:"jsonPropertyName,omitempty"`
	Positional          bool `json:"positional,omitempty"`
	WithInit            bool `json:"withInit,omitempty"`
}

type MyBatisOptions struct {
	AllCrud bool `json:"allCrud"`
	Select  bool `json:"select"`
	Insert  bool `json:"insert"`
	Update  bool `json:"update"`
	Delete  bool `json:"delete"`
}

type TypeScriptOptions struct {
	ExportAllTypes     bool `json:"exportAllTypes,omitempty"`
	ReadonlyProperties bool `json:"readonlyProperties,omitempty"`
	OptionalProperties bool `json:"optionalProperties,omitempty"`
	StrictNullChecks   bool `json:"strictNullChecks,omitempty"`
	Comments           bool `json:"comments,omitempty"`
	JSDocComments      bool `json:"jsDocComments,omitempty"`
	PartialType        bool `json:"partialType,omitempty"`
	ReadonlyType       bool `json:"readonlyType,omitempty"`
	ExtraSpacing       bool `json:"extraSpacing,omitempty"`
}

type ZodOptions struct {
	ExportAllTypes bool `json:"exportAllTypes,omitempty"`
	AllOptional    bool `json:"allOptional,omitempty"`
	Comments       bool `json:"comments,omitempty"`
	Nullable       bool `json:"nullable,omitempty"`
	Nullish        bool `json:"nullish,omitempty"`
	MaxValue       bool `json:"maxValue,omitempty"`
	Trim           bool `json:"trim,omitempty"`
}

type DatabaseConnectionHealth struct {
	ConnectionID  int        `json:"connectionId"`
	Name          string     `json:"name"`
	Status        string     `json:"status"`
	LastCheckedAt *time.Time `json:"lastCheckedAt,omitempty"`
}

const (
	Healthy   HealthStatus = "healthy"
	Degraded  HealthStatus = "degraded"
	Unhealthy HealthStatus = "unhealthy"
)

type DatabaseHealth struct {
	Connected bool  `json:"connected"`
	LatencyMS int64 `json:"latency"`
}

type HealthResponse struct {
	Status   HealthStatus   `json:"status"`
	Version  string         `json:"version"`
	Uptime   int64          `json:"uptime"`
	Database DatabaseHealth `json:"database"`
}
