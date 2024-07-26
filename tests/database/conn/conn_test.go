package conn_test

import (
	"database/sql"
	"testing"

	"github.com/A-Victory/blog/database"
	"github.com/A-Victory/blog/database/conn"
	_ "github.com/go-sql-driver/mysql"
)

const (
	testDBConfig     = "root:password@tcp(127.0.0.1:3306)/"
	testDBName       = "testdb"
	testDBConnection = testDBConfig + testDBName
)

// setupTestDB sets up a connection to the MySQL server, creates the test database, and connects to it.
func setupTestDB(t *testing.T) *database.DBconn {
	dbConn := database.NewDBConn(testDBConfig, testDBName)
	if dbConn == nil || dbConn.DB == nil {
		t.Fatalf("Failed to create a new DB connection")
	}

	// Ensure tables are clean
	cleanupTestDB(dbConn.DB, t, testDBName)
	return dbConn
}

// cleanupTestDB cleans up the test database by dropping tables and the database itself.
func cleanupTestDB(db *sql.DB, t *testing.T, dbName string) {
	// Use the specific database context
	/*
		_, err := db.Exec("USE " + dbName)
		if err != nil {
			t.Fatalf("Failed to switch to test database: %v", err)
		}
	*/

	_, err := db.Exec("DROP TABLE IF EXISTS comments, posts, users")
	if err != nil {
		t.Fatalf("Failed to clean up test database tables: %v", err)
	}

	// Switch back to the root context to drop the database
	_, err = db.Exec("USE mysql")
	if err != nil {
		t.Fatalf("Failed to switch to root database: %v", err)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		t.Fatalf("Failed to drop test database: %v", err)
	}
}

// TestNewConn tests the NewConn function to ensure it creates a DB struct with a non-nil database connection.
func TestNewConn(t *testing.T) {
	// Setup test database
	dbConn := setupTestDB(t)

	// Create a new DB instance using NewConn
	db := conn.NewConn(dbConn)

	// Check if the conn field is correctly set
	if db.Conn == nil {
		t.Fatal("Expected conn to be non-nil")
	}
	if db.Conn != dbConn {
		t.Fatalf("Expected conn to be set to the provided database connection, got %v", db.Conn)
	}

	// Perform an additional check to ensure the connection is valid
	err := db.Conn.DB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Cleanup test database
	cleanupTestDB(dbConn.DB, t, testDBName)
}
