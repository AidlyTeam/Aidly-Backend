package dbadapter

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
)

// DatabaseType defines a custom string type for database types
type DatabaseType string

// Constants for supported database types
const (
	MYSQL      DatabaseType = "mysql"
	POSTGRESQL DatabaseType = "postgres"
)

// idbAdapter interface
type idbAdapter interface {
	Connect() (*sqlx.DB, error)                // Method to establish a database connection
	ConnectAndMigrateGoose() (*sqlx.DB, error) // ConnectAndMigrateGoose connects to the database and applies migrations using Goose
	Migrate(db *sql.DB) error                  // Migrate applies migrations to the given database using Goose
}

// dbAdapter struct
type dbAdapter struct {
	dbType        DatabaseType // Type of the database
	driver        string       // Database driver (e.g., "mysql", "postgres")
	connectionStr string       // Connection string for the database
	migrationPath string       // Migration path
}

// Init function creates a new dbAdapter instance
func Init(dbType DatabaseType, driver, connectionStr, migrationPath string) (idbAdapter, error) {
	switch dbType {
	case MYSQL, POSTGRESQL:
		return &dbAdapter{dbType: dbType, driver: driver, connectionStr: connectionStr, migrationPath: migrationPath}, nil
	default:
		return nil, errors.New("unsupported database type")
	}
}

// Connect method establishes a connection to the database
func (c *dbAdapter) Connect() (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	// Call the appropriate connection method based on the database type
	switch c.dbType {
	case MYSQL:
		db, err = c.newMySQL()
	case POSTGRESQL:
		db, err = c.newPostgres()
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *dbAdapter) ConnectAndMigrateGoose() (conn *sqlx.DB, err error) {
	conn, err = c.Connect()
	if err != nil {
		return nil, err
	}

	if err = c.Migrate(conn.DB); err != nil {
		return
	}

	return conn, nil
}

// newMySQL method opens a connection for MySQL
func (c *dbAdapter) newMySQL() (conn *sqlx.DB, err error) {
	conn, err = sqlx.Connect(c.driver, c.connectionStr) // Establish MySQL connection
	if err != nil {
		return
	}

	// Test the database connection
	if err = conn.Ping(); err != nil {
		return
	}

	time.Sleep(5 * time.Second)
	return
}

// newPostgres method opens a connection for PostgreSQL
func (c *dbAdapter) newPostgres() (conn *sqlx.DB, err error) {
	conn, err = sqlx.Connect(c.driver, c.connectionStr) // Establish PostgreSQL connection
	if err != nil {
		return
	}

	// Test the database connection
	if err = conn.Ping(); err != nil && err.Error() != "pq: database system is starting up" {
		return
	}

	time.Sleep(5 * time.Second)
	return
}

// Migrate method performs database migrations using Goose
func (c *dbAdapter) Migrate(db *sql.DB) error {
	err := goose.SetDialect(string(c.dbType))
	if err != nil {
		return err
	}
	if err := goose.Up(db, c.migrationPath); err != nil {
		return err
	}

	return nil
}
