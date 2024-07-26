package conn

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/A-Victory/blog/models"
)

func (db *DB) CreatePost(data models.Post) (int, error) {

	data.CreatedAt = time.Now().Local().Format("2006-01-02 15:04:05")
	data.UpdatedAt = time.Now().UTC().Format("2006-01-02 15:04:05")
	query := "INSERT INTO Posts (title, content, authorId, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?)"
	result, err := db.Conn.DB.Exec(query, data.Title, data.Content, data.AuthorID, data.CreatedAt, data.UpdatedAt)
	if err != nil {
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(postID), nil
}

func (db *DB) DeletePost(postID int) (int, error) {
	query := "DELETE FROM posts WHERE id = ?"
	result, err := db.Conn.DB.Exec(query, postID)
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

func (db *DB) UpdatePost(post models.Post) (int, error) {
	query := "UPDATE posts SET "
	params := []interface{}{}

	if post.Title != "" {
		query += "title = ?, "
		params = append(params, post.Title)
	}

	if post.Content != "" {
		query += "content = ?, "
		params = append(params, post.Content)
	}

	if len(params) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}

	updatedAt := time.Now().UTC().Format("2006-01-02 15:04:05")

	query += "updatedAt = ? WHERE id = ? AND authorId = ?"
	params = append(params, updatedAt, post.ID, post.AuthorID)

	log.Printf("%+v", params)

	result, err := db.Conn.DB.Exec(query, params...)
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

	//id, _ := result.LastInsertId()

	return int(rowsAffected), nil
}

func (db *DB) GetPosts(limit, offset *int, searchTerm *string) ([]models.Post, error) {
	query := `SELECT id, title, content, authorId, createdAt, updatedAt FROM posts WHERE 1=1`

	params := []interface{}{}

	if searchTerm != nil && *searchTerm != "" {
		query += " AND (title LIKE ? OR content LIKE ?)"
		searchValue := "%" + strings.TrimSpace(*searchTerm) + "%"
		params = append(params, searchValue, searchValue)
	}

	if limit != nil && offset != nil {
		query += " LIMIT ? OFFSET ?"
		params = append(params, *limit, *offset)
	}

	stmt, err := db.Conn.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (db *DB) GetPostByID(postID int) (models.Post, error) {
	query := "SELECT id, title, content, authorId, createdAt, updatedAt FROM posts WHERE id = ?"
	row := db.Conn.DB.QueryRow(query, postID)

	var post models.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// return models.Post{}, fmt.Errorf("no post found with ID %d", postID)
			return models.Post{}, nil
		}
		return models.Post{}, err
	}

	return post, nil
}
