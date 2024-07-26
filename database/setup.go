package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DBconn struct {
	DB *sql.DB
}

func NewDBConn(config, dbName string) *DBconn {

	db, err := sql.Open("mysql", config)
	if err != nil {
		log.Fatal("unable to open database: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("unable to connect to database: ", err)
	}

	// query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	query := "CREATE DATABASE IF NOT EXISTS " + dbName
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("failed to create database %s: %v", dbName, err)
	}

	_, err = db.Exec("USE " + dbName)
	if err != nil {
		log.Fatalf("failed to use database %s: %v", dbName, err)
	}

	return &DBconn{
		DB: db,
	}
}

func (dbConn *DBconn) Initialize() error {

	/*
		tx, err := dbConn.DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				panic(err)
			}
		}()
	*/

	createUserTable := `
	CREATE TABLE IF NOT EXISTS Users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);`

	// Create the Post table
	createPostTable := `
	CREATE TABLE IF NOT EXISTS Posts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		authorId INT NOT NULL,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (authorId) REFERENCES Users(id) ON DELETE CASCADE
	);`

	// Create the Comment table
	createCommentTable := `
	CREATE TABLE IF NOT EXISTS Comments (
		id INT AUTO_INCREMENT PRIMARY KEY,
		postId INT NOT NULL,
		authorId INT NOT NULL,
		content TEXT NOT NULL,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (postId) REFERENCES Posts(id) ON DELETE CASCADE,
		FOREIGN KEY (authorId) REFERENCES Users(id) ON DELETE CASCADE
	);`

	// Execute the table creation statements
	_, err := dbConn.DB.Exec(createUserTable)
	if err != nil {
		return err
	}

	_, err = dbConn.DB.Exec(createPostTable)
	if err != nil {
		return err
	}

	_, err = dbConn.DB.Exec(createCommentTable)
	if err != nil {
		return err
	}

	log.Println("Tables created successfully!")
	return nil
}
