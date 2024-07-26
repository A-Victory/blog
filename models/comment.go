package models

type Comment struct {
	ID        int    `json:"id"`
	Postid    int    `json:"postId"`
	AuthorID  int    `json:"authorId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
