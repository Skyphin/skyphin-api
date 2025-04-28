package models

type VoteType string

const (
	Upvote   VoteType = "upvote"
	Downvote VoteType = "downvote"
)

type CommentVote struct {
	CommentID string   `json:"comment_id" gorm:"primaryKey"`
	UserID    string   `json:"user_id" gorm:"primaryKey"`
	VoteType  VoteType `json:"vote_type"`
	CreatedAt string   `json:"created_at"`
}

type URLVote struct {
	UserID    string   `json:"user_id" gorm:"primaryKey"`
	URL       string   `json:"url" gorm:"primaryKey"`
	VoteType  VoteType `json:"vote_type"`
	CreatedAt string   `json:"created_at"`
}

type URLVoteResult struct {
	URL           string `json:"url"`
	UpvoteCount   int    `json:"unvote_count"`
	DownvoteCount int    `json:"downvote_count"`
}
