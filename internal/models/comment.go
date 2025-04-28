package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	URL             string         `json:"url" gorm:"type:text;not null;index"`
	ParentCommentID *string        `json:"parent_comment_id" gorm:"type:uuid;index"`
	Content         string         `json:"content" gorm:"type:text;not null"`
	AuthorID        string         `json:"author_id" gorm:"type:uuid;not null;index"`
	CreatedAt       time.Time      `json:"created_at" gorm:"type:timestamptz;not null;default:now()"`
	UpvoteCount     int            `json:"upvote_count" gorm:"type:integer;not null;default:0"`
	DownvoteCount   int            `json:"downvote_count" gorm:"type:integer;not null;default:0"`
	ReplyCount      int            `json:"reply_count" gorm:"type:integer;not null;default:0"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type AddCommentInput struct {
	Content  string `json:"content" validate:"required,min=1,max=5000"`
	URL      string `json:"url" validate:"required,url"`
	AuthorID string `json:"author_id" validate:"required,uuid4"`
}

type AddReplyInput struct {
	Content         string `json:"content" validate:"required,min=1,max=5000"`
	URL             string `json:"url" validate:"required,url"`
	ParentCommentID string `json:"parent_comment_id" validate:"required,uuid4"`
	AuthorID        string `json:"author_id" validate:"required,uuid4"`
}

type PageInfo struct {
	HasNextPage bool    `json:"has_next_page"`
	EndCursor   *string `json:"end_cursor"`
}

type CommentConnection struct {
	Edges      []*CommentEdge `json:"edges"`
	PageInfo   *PageInfo      `json:"page_info"`
	TotalCount int            `json:"total_count"`
}

type CommentEdge struct {
	Node   *Comment `json:"node"`
	Cursor string   `json:"cursor"`
}
