package conn

import (
	"time"

	"github.com/A-Victory/blog/models"
)

func (db *DB) AddComment(data models.Comment) (int, error) {

	data.CreatedAt = time.Now().Local().Format("2006-01-02 15:04:05")
	data.UpdatedAt = time.Now().Local().Format("2006-01-02 15:04:05")
	query := "INSERT INTO comments (postId, authorId, content, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?)"

	result, err := db.Conn.DB.Exec(query, data.Postid, data.AuthorID, data.Content, data.CreatedAt, data.UpdatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, nil
	}

	comment_id, _ := result.LastInsertId()

	return int(comment_id), nil
}

func (db *DB) DeleteComment(commentID int) (int, error) {
	query := "DELETE FROM comments WHERE id = ?"
	result, err := db.Conn.DB.Exec(query, commentID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, nil // No rows affected, indicating no post with the given ID was found
	}

	return int(rowsAffected), nil
}

func (db *DB) EditComment(comment models.Comment) (int, error) {
	query := "UPDATE comments SET content = ?, updatedAt = ? WHERE id = ? AND authorId = ? AND postId = ?"
	result, err := db.Conn.DB.Exec(query, comment.Content, time.Now().UTC(), comment.ID, comment.AuthorID, comment.Postid)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, nil // No rows affected, indicating no comment with the given criteria was found
	}

	return int(rowsAffected), nil
}

func (db *DB) GetComments(postID, limit, offset int) ([]models.Comment, error) {
	query := "SELECT id, postId, authorId, content, createdAt, updatedAt FROM comments WHERE postId = ? LIMIT ? OFFSET ?"
	rows, err := db.Conn.DB.Query(query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Postid, &comment.AuthorID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
