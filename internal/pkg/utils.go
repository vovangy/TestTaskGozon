package pkg

import (
	"crypto/sha1"
	"encoding/hex"
	"myHabr/internal/models"
	"strconv"
)

func GenerateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func ConvertToNestedComments(tree *models.CommentTree) []*models.Comment {
	rootComment := &models.Comment{
		ID:      strconv.FormatInt(tree.Comment.CommentId, 10),
		Content: tree.Comment.Content,
		Author:  nil,
		Replies: []*models.Comment{},
	}

	for _, replyTree := range tree.Replies {
		rootComment.Replies = append(rootComment.Replies, ConvertToNestedComments(replyTree)...)
	}

	return []*models.Comment{rootComment}
}
