package models

type User struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Posts    []*Post    `json:"posts"`
	Comments []*Comment `json:"comments"`
}
