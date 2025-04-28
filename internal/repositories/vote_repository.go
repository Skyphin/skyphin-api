package repositories

import (
	"context"
	"database/sql"
	"errors"
	"skyphin-api/internal/models"
)

type VoteRepository struct {
	db *sql.DB
}

func NewVoteRepository(db *sql.DB) *VoteRepository {
	return &VoteRepository{db: db}
}

func (r *VoteRepository) GetCommentVote(ctx context.Context, userID, commentID string) (*models.CommentVote, error) {
	query := `
		SELECT user_id, comment_id, vote_type, created_at
		FROM comment_votes
		WHERE user_id = $1 AND comment_id = $2
	`
	row := r.db.QueryRowContext(ctx, query, userID, commentID)

	var vote models.CommentVote
	err := row.Scan(
		&vote.UserID,
		&vote.CommentID,
		&vote.VoteType,
		&vote.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *VoteRepository) CreateCommentVote(ctx context.Context, vote *models.CommentVote) error {
	query := `
		INSERT INTO comment_votes (user_id, comment_id, vote_type, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query,
		vote.UserID,
		vote.CommentID,
		vote.VoteType,
		vote.CreatedAt,
	)
	return err
}

func (r *VoteRepository) UpdateCommentVote(ctx context.Context, userID, commentID string, voteType models.VoteType) error {
	query := `
		UPDATE comment_votes
		SET vote_type = $1
		WHERE user_id = $2 AND comment_id = $3
	`
	_, err := r.db.ExecContext(ctx, query, voteType, userID, commentID)
	return err
}

func (r *VoteRepository) DeleteCommentVote(ctx context.Context, userID, commentID string) error {
	query := `
		DELETE FROM comment_votes
		WHERE user_id = $1 AND comment_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, userID, commentID)
	return err
}

func (r *VoteRepository) GetURLVote(ctx context.Context, userID, url string) (*models.URLVote, error) {
	query := `
		SELECT user_id, url, vote_type, created_at
		FROM url_votes
		WHERE user_id = $1 AND url = $2
	`
	row := r.db.QueryRowContext(ctx, query, userID, url)

	var vote models.URLVote
	err := row.Scan(
		&vote.UserID,
		&vote.URL,
		&vote.VoteType,
		&vote.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

func (r *VoteRepository) CreateURLVote(ctx context.Context, vote *models.URLVote) error {
	query := `
		INSERT INTO url_votes (user_id, url, vote_type, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query,
		vote.UserID,
		vote.URL,
		vote.VoteType,
		vote.CreatedAt,
	)
	return err
}

func (r *VoteRepository) UpdateURLVote(ctx context.Context, userID, url string, voteType models.VoteType) error {
	query := `
		UPDATE url_votes
		SET vote_type = $1
		WHERE user_id = $2 AND url = $3
	`
	_, err := r.db.ExecContext(ctx, query, voteType, userID, url)
	return err
}

func (r *VoteRepository) DeleteURLVote(ctx context.Context, userID, url string) error {
	query := `
		DELETE FROM url_votes
		WHERE user_id = $1 AND url = $2
	`
	_, err := r.db.ExecContext(ctx, query, userID, url)
	return err
}

func (r *VoteRepository) GetURLVoteCounts(ctx context.Context, url string) (int, int, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN vote_type = 'upvote' THEN 1 END) as upvotes,
			COUNT(CASE WHEN vote_type = 'downvote' THEN 1 END) as downvotes
		FROM url_votes
		WHERE url = $1
	`
	var upvotes, downvotes int
	err := r.db.QueryRowContext(ctx, query, url).Scan(&upvotes, &downvotes)
	return upvotes, downvotes, err
}

func (r *VoteRepository) GetCommentVoteCounts(ctx context.Context, commentID string) (int, int, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN vote_type = 'upvote' THEN 1 END) as upvotes,
			COUNT(CASE WHEN vote_type = 'downvote' THEN 1 END) as downvotes
		FROM comment_votes
		WHERE comment_id = $1
	`
	var upvotes, downvotes int
	err := r.db.QueryRowContext(ctx, query, commentID).Scan(&upvotes, &downvotes)
	return upvotes, downvotes, err
}
