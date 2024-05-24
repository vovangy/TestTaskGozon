package models

type Comment struct {
	ID            string     `json:"id"`
	Content       string     `json:"content"`
	Author        *User      `json:"author"`
	Post          *Post      `json:"post,omitempty"`
	ParentComment *Comment   `json:"parentComment,omitempty"`
	Replies       []*Comment `json:"replies"`
}
