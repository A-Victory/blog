package conn_test

import (
	"testing"
	"time"

	"github.com/A-Victory/blog/database/conn"
	"github.com/A-Victory/blog/models"
	_ "github.com/go-sql-driver/mysql"
)

// TestCommentFunctions tests AddComment, DeleteComment, EditComment, and GetComments functions.
func TestCommentFunctions(t *testing.T) {
	// Setup test database
	dbConn := setupTestDB(t)
	defer cleanupTestDB(dbConn.DB, t, testDBName)

	// Initialize database
	err := dbConn.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	db := conn.NewConn(dbConn)

	// Define a user to save (for author_id in comments)
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	userID, err := db.SaveUser(user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Define a post to create (for post_id in comments)
	post := models.Post{
		Title:     "Test Post",
		Content:   "This is a test post.",
		AuthorID:  userID,
		CreatedAt: time.Now().Local().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Local().Format("2006-01-02 15:04:05"),
	}
	postID, err := db.CreatePost(post)
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Define a comment to add
	comment := models.Comment{
		Postid:    postID,
		AuthorID:  userID,
		Content:   "This is a test comment.",
		CreatedAt: time.Now().Local().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Local().Format("2006-01-02 15:04:05"),
	}

	// Test AddComment
	commentID, err := db.AddComment(comment)
	if err != nil {
		t.Fatalf("Failed to add comment: %v", err)
	}
	if commentID <= 0 {
		t.Fatalf("Invalid comment ID returned: %d", commentID)
	}

	// Test GetComments
	comments, err := db.GetComments(postID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get comments: %v", err)
	}
	if len(comments) == 0 {
		t.Fatalf("Expected at least one comment, got %d", len(comments))
	}

	// Test GetCommentByID
	savedComment, err := db.GetCommentByID(commentID)
	if err != nil {
		t.Fatalf("Failed to get comment by ID: %v", err)
	}
	if savedComment.ID != commentID {
		t.Fatalf("Expected comment ID %d, got %d", commentID, savedComment.ID)
	}

	// Test EditComment
	comment.ID = commentID
	comment.Content = "Updated test comment."
	id, err := db.EditComment(comment)
	if err != nil {
		t.Fatalf("Failed to edit comment: %v", err)
	}
	if id <= 0 {
		t.Fatalf("Expected to update comment with ID %d, but found %d rows", postID, id)
	}

	// Verify the comment is updated
	updatedComment, err := db.GetCommentByID(commentID)
	if err != nil {
		t.Fatalf("Failed to get comment by ID: %v", err)
	}
	if updatedComment.Content != "Updated test comment." {
		t.Fatalf("Expected updated comment content, got %s", updatedComment.Content)
	}

	// Test DeleteComment
	deletedCommentID, err := db.DeleteComment(commentID)
	if err != nil {
		t.Fatalf("Failed to delete comment: %v", err)
	}
	if deletedCommentID != commentID {
		t.Fatalf("Expected deleted comment ID %d, got %d", commentID, deletedCommentID)
	}

	// Verify the comment is deleted
	comments, err = db.GetComments(postID, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get comments: %v", err)
	}
	for _, c := range comments {
		if c.ID == commentID {
			t.Fatal("Expected comment to be deleted, but it still exists")
		}
	}

	// Verify the comment is not retrievable by ID
	deletedComment, err := db.GetCommentByID(commentID)
	if err != nil {
		t.Fatalf("Failed to get comment by ID: %v", err)
	}
	if deletedComment.ID != 0 {
		t.Fatal("Expected no comment to be found, but found one")
	}
}
