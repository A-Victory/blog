package conn_test

import (
	"testing"

	"github.com/A-Victory/blog/database/conn"
	"github.com/A-Victory/blog/models"
	_ "github.com/go-sql-driver/mysql"
)

// TestUserFunctions tests SaveUser and GetUser functions for various cases.
func TestUserFunctions(t *testing.T) {
	// Setup test database
	dbConn := setupTestDB(t)

	// Create a new DB instance using the connection
	err := dbConn.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	db := conn.NewConn(dbConn)

	// Define a user to save
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Test SaveUser
	id, err := db.SaveUser(user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}
	if id <= 0 {
		t.Fatalf("Invalid user ID returned: %d", id)
	}

	// Test GetUser by ID
	retrievedUser, err := db.GetUser("id", id)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}
	if retrievedUser.ID != id {
		t.Fatalf("Expected user ID %d, got %d", id, retrievedUser.ID)
	}

	// Test GetUser by Username
	retrievedUserByUsername, err := db.GetUser("username", user.Username)
	if err != nil {
		t.Fatalf("Failed to get user by username: %v", err)
	}
	if retrievedUserByUsername.Username != user.Username {
		t.Fatalf("Expected username %s, got %s", user.Username, retrievedUserByUsername.Username)
	}

	// Test GetUser by Email
	retrievedUserByEmail, err := db.GetUser("email", user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}
	if retrievedUserByEmail.Email != user.Email {
		t.Fatalf("Expected email %s, got %s", user.Email, retrievedUserByEmail.Email)
	}

	// Test GetUser with invalid identifier type
	_, err = db.GetUser("invalid", "somevalue")
	if err == nil {
		t.Fatal("Expected error for invalid identifier type, got nil")
	}
	if err.Error() != "invalid identifier type" {
		t.Fatalf("Expected error 'invalid identifier type', got %v", err)
	}

	cleanupTestDB(dbConn.DB, t, testDBName)
}
