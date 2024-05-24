package graph

import "myHabr/internal/models"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

var users = []*models.User{{ID: "1", Username: "Vova", Posts: nil, Comments: nil}}
var posts = []*models.Post{{ID: "1", Title: "Aba", Content: "CHTOTO", Author: users[0], Comments: nil}}

type Resolver struct{}
