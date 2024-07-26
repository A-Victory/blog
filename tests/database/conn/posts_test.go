package conn_test

import (
	"testing"

	"github.com/A-Victory/blog/database/conn"
	"github.com/A-Victory/blog/models"
	_ "github.com/go-sql-driver/mysql"
)

// TestPostFunctions tests CreatePost, DeletePost, UpdatePost, GetPosts, and GetPostByID functions.
func TestPostFunctions(t *testing.T) {
	// Setup test database
	dbConn := setupTestDB(t)
	defer cleanupTestDB(dbConn.DB, t, testDBName)

	// Initialize database
	err := dbConn.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	db := conn.NewConn(dbConn)

	// Define a user to save (for author_id in posts)
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	userID, err := db.SaveUser(user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Define a post to create
	post := models.Post{
		Title:    "Test Post",
		Content:  "This is a test post.",
		AuthorID: userID,
	}

	// Test CreatePost
	postID, err := db.CreatePost(post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}
	if postID <= 0 {
		t.Fatalf("Invalid post ID returned: %d", postID)
	}

	// Test GetPostByID
	retrievedPost, err := db.GetPostByID(postID)
	if err != nil {
		t.Fatalf("Failed to get post by ID: %v", err)
	}
	if retrievedPost.ID != postID {
		t.Fatalf("Expected post ID %d, got %d", postID, retrievedPost.ID)
	}

	// Test UpdatePost
	post.ID = postID
	post.Title = "Updated Test Post"
	post.Content = "This is an updated test post."
	post.AuthorID = userID
	updatedPostID, err := db.UpdatePost(post)
	if err != nil {
		t.Fatalf("Failed to update post: %v", err)
	}
	if updatedPostID <= 0 {
		t.Fatalf("Expected to update post ID %d, but found %d rows", postID, updatedPostID)
	}

	// Test GetPosts with pagination and search
	limit := 10
	offset := 0
	searchTerm := "Updated"
	posts, err := db.GetPosts(&limit, &offset, &searchTerm)
	if err != nil {
		t.Fatalf("Failed to get posts with pagination and search: %v", err)
	}
	if len(posts) == 0 {
		t.Fatalf("Expected at least one post, got %d", len(posts))
	}

	// Test DeletePost
	deletedPostID, err := db.DeletePost(postID)
	if err != nil {
		t.Fatalf("Failed to delete post: %v", err)
	}
	if deletedPostID != postID {
		t.Fatalf("Expected deleted post ID %d, got %d", postID, deletedPostID)
	}

	// Verify the post is deleted
	_, err = db.GetPostByID(postID)
	if err != nil {
		t.Fatalf("Expected no error for non-existent post ID, got: %v", err)
	}
}
