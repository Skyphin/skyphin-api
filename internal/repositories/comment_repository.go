package repositories

import (
	"context"
	"database/sql"
	"errors"
	"skyphin-api/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (id, url, parent_comment_id, content, author_id, created_at, upvote_count, downvote_count, reply_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		comment.ID,
		comment.URL,
		comment.ParentCommentID,
		comment.Content,
		comment.AuthorID,
		comment.CreatedAt,
		comment.UpvoteCount,
		comment.DownvoteCount,
		comment.ReplyCount,
	)
	return err
}

func (r *CommentRepository) GetByID(ctx context.Context, id string) (*models.Comment, error) {
	query := `
		SELECT id, url, parent_comment_id, content, author_id, created_at, upvote_count, downvote_count, reply_count
		FROM comments
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var comment models.Comment
	err := row.Scan(
		&comment.ID,
		&comment.URL,
		&comment.ParentCommentID,
		&comment.Content,
		&comment.AuthorID,
		&comment.CreatedAt,
		&comment.UpvoteCount,
		&comment.DownvoteCount,
		&comment.ReplyCount,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) GetByURL(ctx context.Context, url string, first int, after *string) (*models.CommentConnection, error) {
	query := `
		SELECT id, url, parent_comment_id, content, author_id, created_at, upvote_count, downvote_count, reply_count
		FROM comments
		WHERE url = $1 AND parent_comment_id IS NULL
	`

	var args []interface{}
	args = append(args, url)

	// Get total count
	countQuery := `SELECT COUNT(*) FROM comments WHERE url = $1 AND parent_comment_id IS NULL`
	var totalCount int
	if err := r.db.QueryRowContext(ctx, countQuery, url).Scan(&totalCount); err != nil {
		return nil, err
	}

	if after != nil {
		query += " AND id > $2"
		args = append(args, *after)
	}

	query += " ORDER BY created_at DESC"

	if first > 0 {
		query += " LIMIT $"
		if after != nil {
			query += "3"
		} else {
			query += "2"
		}
		args = append(args, first)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edges []*models.CommentEdge
	var endCursor *string

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.URL,
			&comment.ParentCommentID,
			&comment.Content,
			&comment.AuthorID,
			&comment.CreatedAt,
			&comment.UpvoteCount,
			&comment.DownvoteCount,
			&comment.ReplyCount,
		)
		if err != nil {
			return nil, err
		}

		edges = append(edges, &models.CommentEdge{
			Node:   &comment,
			Cursor: comment.ID,
		})
		endCursor = &comment.ID
	}

	hasNextPage := first > 0 && len(edges) == first

	return &models.CommentConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			EndCursor:   endCursor,
		},
		TotalCount: totalCount,
	}, nil
}

func (r *CommentRepository) GetReplies(ctx context.Context, parentID string, first int, after *string) (*models.CommentConnection, error) {
	// Base query
	query := `
		SELECT id, url, parent_comment_id, content, author_id, created_at, upvote_count, downvote_count, reply_count
		FROM comments
		WHERE parent_comment_id = $1
	`

	var args []interface{}
	args = append(args, parentID)

	if after != nil {
		query += " AND id > $2"
		args = append(args, *after)
	}

	// Get total count for pagination
	countQuery := `SELECT COUNT(*) FROM comments WHERE parent_comment_id = $1`
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, parentID).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Add ordering and limit
	query += " ORDER BY created_at ASC"

	if first > 0 {
		query += " LIMIT $"
		if after != nil {
			query += "3"
		} else {
			query += "2"
		}
		args = append(args, first)
	}

	// Execute query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results
	var edges []*models.CommentEdge
	var endCursor *string

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.URL,
			&comment.ParentCommentID,
			&comment.Content,
			&comment.AuthorID,
			&comment.CreatedAt,
			&comment.UpvoteCount,
			&comment.DownvoteCount,
			&comment.ReplyCount,
		)
		if err != nil {
			return nil, err
		}

		edges = append(edges, &models.CommentEdge{
			Node:   &comment,
			Cursor: comment.ID,
		})
		endCursor = &comment.ID
	}

	// Determine if there are more results
	hasNextPage := false
	if first > 0 && len(edges) == first {
		hasNextPage = true
	}

	return &models.CommentConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			EndCursor:   endCursor,
		},
		TotalCount: totalCount,
	}, nil
}

func (r *CommentRepository) IncrementReplyCount(ctx context.Context, parentID string) error {
	query := `
		UPDATE comments
		SET reply_count = reply_count + 1
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, parentID)
	return err
}

func (r *CommentRepository) UpdateVoteCount(ctx context.Context, commentID string, upvoteDelta, downvoteDelta int) error {
	query := `
		UPDATE comments
		SET upvote_count = upvote_count + $1,
			downvote_count = downvote_count + $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, upvoteDelta, downvoteDelta, commentID)
	return err
}

func (r *CommentRepository) CountByURL(ctx context.Context, url string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM comments
		WHERE url = $1
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, url).Scan(&count)
	return count, err
}
