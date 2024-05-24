package models

type Post struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Content  string     `json:"content"`
	Author   *User      `json:"author"`
	Comments []*Comment `json:"comments"`
}
