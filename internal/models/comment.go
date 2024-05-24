package models

type Comment struct {
	ID            string     `json:"id"`
	Content       string     `json:"content"`
	Author        *User      `json:"author"`
	Post          *Post      `json:"post,omitempty"`
	ParentComment *Comment   `json:"parentComment,omitempty"`
	Replies       []*Comment `json:"replies"`
}

type CommentCreateData struct {
	UserId          int64
	PostId          int64
	CommentParentId int64
	Content         string
}

type CommentCreateResponse struct {
	CommentId       int64
	UserId          int64
	PostId          int64
	CommentParentId int64
	Content         string
}

type CommentsBlockRequest struct {
	UserId int64
	PostId int64
}

type CommentTree struct {
	Comment CommentCreateResponse
	Replies []CommentTree
}
