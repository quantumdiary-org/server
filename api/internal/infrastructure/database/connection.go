package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBType represents the type of database
type DBType string

const (
	PostgreSQL DBType = "postgres"
	MariaDB    DBType = "mariadb"
	MySQL      DBType = "mysql"
	SQLite     DBType = "sqlite"
)

// ConnectionManager manages connections to different database types
type ConnectionManager struct {
	config DatabaseConfig
}

// DatabaseConfig holds configuration for database connections
type DatabaseConfig struct {
	Type     string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
	URL      string
	SQLitePath string
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(config DatabaseConfig) *ConnectionManager {
	return &ConnectionManager{config: config}
}

// Connect creates a connection to the database based on the configured type
func (cm *ConnectionManager) Connect() (*gorm.DB, error) {
	dbType := DBType(cm.config.Type)

	switch dbType {
	case PostgreSQL:
		return cm.connectPostgreSQL()
	case MariaDB, MySQL:
		return cm.connectMariaDB()
	case SQLite:
		return cm.connectSQLite()
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// connectPostgreSQL creates a connection to PostgreSQL
func (cm *ConnectionManager) connectPostgreSQL() (*gorm.DB, error) {
	var dsn string
	if cm.config.URL != "" {
		dsn = cm.config.URL
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cm.config.Host, cm.config.User, cm.config.Password,
			cm.config.Name, cm.config.Port, cm.config.SSLMode)
	}

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// connectMariaDB creates a connection to MariaDB/MySQL
func (cm *ConnectionManager) connectMariaDB() (*gorm.DB, error) {
	var dsn string
	if cm.config.URL != "" {
		dsn = cm.config.URL
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cm.config.User, cm.config.Password, cm.config.Host, cm.config.Port, cm.config.Name)
	}

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// connectSQLite creates a connection to SQLite
func (cm *ConnectionManager) connectSQLite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(cm.config.SQLitePath), &gorm.Config{})
}