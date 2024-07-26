package setup_test

import (
	"database/sql"
	"testing"

	"github.com/A-Victory/blog/database"
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
	//cleanupTestDB(dbConn.DB, t)
	return dbConn
}

// TestNewDBConn tests the NewDBConn method for establishing a database connection.
func TestNewDBConn(t *testing.T) {
	dbConn := setupTestDB(t)
	defer cleanupTestDB(dbConn.DB, t)

	err := dbConn.DB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

// cleanupTestDB cleans up the test database by dropping tables and the database itself.
func cleanupTestDB(db *sql.DB, t *testing.T) {
	_, err := db.Exec("DROP TABLE IF EXISTS comments, posts, users")
	if err != nil {
		t.Fatalf("Failed to clean up test database tables: %v", err)
	}
}

// TestInitialize tests the Initialize method to create tables in the database.
func TestInitialize(t *testing.T) {
	dbConn := setupTestDB(t)

	err := dbConn.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Query to list all tables in the database
	rows, err := dbConn.DB.Query("SHOW TABLES")
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	defer rows.Close()

	// Expected table names
	tables := map[string]bool{
		"users":    false,
		"posts":    false,
		"comments": false,
	}

	// Iterate over the rows to check if the tables exist
	for rows.Next() {
		var existingTableName string
		if err := rows.Scan(&existingTableName); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		t.Logf("Found table: %s", existingTableName)
		if _, ok := tables[existingTableName]; ok {
			tables[existingTableName] = true
		}
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("Error iterating over rows: %v", err)
	}

	// Check for missing tables
	for tableName, exists := range tables {
		if !exists {
			t.Fatalf("Table %s does not exist", tableName)
		}
	}

	// Cleanup
	cleanupTestDB(dbConn.DB, t)

}

// TestUserConstraints tests constraints related to the User table.
func TestUserConstraints(t *testing.T) {
	dbConn := setupTestDB(t)
	defer cleanupTestDB(dbConn.DB, t)

	err := dbConn.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Insert a user
	_, err = dbConn.DB.Exec(`INSERT INTO Users (username, email, password) VALUES ('testuser', 'test@example.com', 'password')`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Try to insert a duplicate user
	_, err = dbConn.DB.Exec(`INSERT INTO Users (username, email, password) VALUES ('testuser', 'test@example.com', 'password')`)
	if err == nil {
		t.Fatal("Expected error for duplicate user insertion, got none")
	}
}

// TestPostConstraints tests constraints related to the Post table.
func TestPostConstraints(t *testing.T) {
	dbConn := setupTestDB(t)
	defer cleanupTestDB(dbConn.DB, t)

	err := dbConn.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Insert a user
	result, err := dbConn.DB.Exec(`INSERT INTO Users (username, email, password) VALUES ('testuser', 'test@example.com', 'password')`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert id: %v", err)
	}

	// Insert a post
	_, err = dbConn.DB.Exec(`INSERT INTO Posts (title, content, authorId) VALUES ('test title', 'test content', ?)`, userID)
	if err != nil {
		t.Fatalf("Failed to insert post: %v", err)
	}

	// Try to insert a post with a non-existent user
	_, err = dbConn.DB.Exec(`INSERT INTO Posts (title, content, authorId) VALUES ('test title', 'test content', 999999)`)
	if err == nil {
		t.Fatal("Expected error for insertion with non-existent user, got none")
	}
}

// TestCommentConstraints tests constraints related to the Comment table.
func TestCommentConstraints(t *testing.T) {
	dbConn := setupTestDB(t)
	defer cleanupTestDB(dbConn.DB, t)

	err := dbConn.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Insert a user
	result, err := dbConn.DB.Exec(`INSERT INTO Users (username, email, password) VALUES ('testuser', 'test@example.com', 'password')`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert id: %v", err)
	}

	// Insert a post
	postResult, err := dbConn.DB.Exec(`INSERT INTO Posts (title, content, authorId) VALUES ('test title', 'test content', ?)`, userID)
	if err != nil {
		t.Fatalf("Failed to insert post: %v", err)
	}

	postID, err := postResult.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert id: %v", err)
	}

	// Insert a comment
	_, err = dbConn.DB.Exec(`INSERT INTO Comments (postId, authorId, content) VALUES (?, ?, 'test comment')`, postID, userID)
	if err != nil {
		t.Fatalf("Failed to insert comment: %v", err)
	}

	// Try to insert a comment with a non-existent post
	_, err = dbConn.DB.Exec(`INSERT INTO Comments (postId, authorId, content) VALUES (999999, ?, 'test comment')`, userID)
	if err == nil {
		t.Fatal("Expected error for insertion with non-existent post, got none")
	}

	// Try to insert a comment with a non-existent user
	_, err = dbConn.DB.Exec(`INSERT INTO Comments (postId, authorId, content) VALUES (?, 999999, 'test comment')`, postID)
	if err == nil {
		t.Fatal("Expected error for insertion with non-existent user, got none")
	}
}
