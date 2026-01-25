package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


type DBType string

const (
	PostgreSQL DBType = "postgres"
	MariaDB    DBType = "mariadb"
	MySQL      DBType = "mysql"
	SQLite     DBType = "sqlite"
)


type ConnectionManager struct {
	config DatabaseConfig
}


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


func NewConnectionManager(config DatabaseConfig) *ConnectionManager {
	return &ConnectionManager{config: config}
}


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


func (cm *ConnectionManager) connectSQLite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(cm.config.SQLitePath), &gorm.Config{})
}