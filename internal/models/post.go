package models

type Post struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Content  string     `json:"content"`
	Author   *User      `json:"author"`
	Comments []*Comment `json:"comments"`
}

type PostCreateData struct {
	UserId      int64
	IsCommented bool
	Title       string
	Content     string
}

type PostCreateResponse struct {
	ID          int64
	UserId      int64
	IsCommented bool
	Title       string
	Content     string
}

type PostResponse struct {
	Post     PostCreateResponse
	Comments []*CommentTree
}
