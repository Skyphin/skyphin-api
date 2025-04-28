package services

import (
	"context"
	"errors"
	"skyphin-api/internal/models"
	"skyphin-api/internal/repositories"
	"skyphin-api/internal/utils"
	"time"

	"github.com/google/uuid"
)

type CommentService struct {
	commentRepo *repositories.CommentRepository
	voteRepo    *repositories.VoteRepository
	userRepo    *repositories.UserRepository
	profanity   *utils.ProfanityFilter
	urlService  *utils.URLService
	rateLimiter *utils.RateLimiter
}

func NewCommentService(
	commentRepo *repositories.CommentRepository,
	voteRepo *repositories.VoteRepository,
	userRepo *repositories.UserRepository,
	profanity *utils.ProfanityFilter,
	urlService *utils.URLService,
	rateLimiter *utils.RateLimiter,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		voteRepo:    voteRepo,
		userRepo:    userRepo,
		profanity:   profanity,
		urlService:  urlService,
		rateLimiter: rateLimiter,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, input *models.AddCommentInput) (*models.Comment, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, errors.New("unauthorized")
	}

	// Rate limiting
	if err := s.rateLimiter.Check(ctx, userID, "comment"); err != nil {
		return nil, err
	}

	// Validate URL
	normalizedURL, err := s.urlService.Normalize(input.URL)
	if err != nil {
		return nil, errors.New("invalid URL")
	}

	// Check profanity
	if s.profanity.HasProfanity(input.Content) {
		return nil, errors.New("comment contains inappropriate content")
	}

	comment := &models.Comment{
		ID:              uuid.New().String(),
		URL:             normalizedURL,
		ParentCommentID: nil,
		Content:         input.Content,
		AuthorID:        userID,
		CreatedAt:       time.Now(),
		UpvoteCount:     0,
		DownvoteCount:   0,
		ReplyCount:      0,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) CreateReply(ctx context.Context, input *models.AddReplyInput) (*models.Comment, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, errors.New("unauthorized")
	}

	// Rate limiting
	if err := s.rateLimiter.Check(ctx, userID, "comment"); err != nil {
		return nil, err
	}

	// Validate parent comment exists
	parent, err := s.commentRepo.GetByID(ctx, input.ParentCommentID)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, errors.New("parent comment not found")
	}

	// Validate URL matches parent
	normalizedURL, err := s.urlService.Normalize(input.URL)
	if err != nil {
		return nil, errors.New("invalid URL")
	}
	if normalizedURL != parent.URL {
		return nil, errors.New("URL mismatch with parent comment")
	}

	// Check profanity
	if s.profanity.HasProfanity(input.Content) {
		return nil, errors.New("comment contains inappropriate content")
	}

	comment := &models.Comment{
		ID:              uuid.New().String(),
		URL:             normalizedURL,
		ParentCommentID: &input.ParentCommentID,
		Content:         input.Content,
		AuthorID:        userID,
		CreatedAt:       time.Now(),
		UpvoteCount:     0,
		DownvoteCount:   0,
		ReplyCount:      0,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Increment parent's reply count
	if err := s.commentRepo.IncrementReplyCount(ctx, input.ParentCommentID); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) GetComment(ctx context.Context, id string) (*models.Comment, error) {
	return s.commentRepo.GetByID(ctx, id)
}

func (s *CommentService) GetCommentsByURL(ctx context.Context, url string, first int, after *string) (*models.CommentConnection, error) {
	normalizedURL, err := s.urlService.Normalize(url)
	if err != nil {
		return nil, errors.New("invalid URL")
	}

	return s.commentRepo.GetByURL(ctx, normalizedURL, first, after)
}

func (s *CommentService) GetReplies(ctx context.Context, parentID string, first int, after *string) (*models.CommentConnection, error) {
	return s.commentRepo.GetReplies(ctx, parentID, first, after)
}

func (s *CommentService) CountCommentsByURL(ctx context.Context, url string) (int, error) {
	normalizedURL, err := s.urlService.Normalize(url)
	if err != nil {
		return 0, errors.New("invalid URL")
	}

	return s.commentRepo.CountByURL(ctx, normalizedURL)
}

func (s *CommentService) UpvoteComment(ctx context.Context, commentID string) (*models.Comment, error) {
	return s.voteComment(ctx, commentID, models.Upvote)
}

func (s *CommentService) DownvoteComment(ctx context.Context, commentID string) (*models.Comment, error) {
	return s.voteComment(ctx, commentID, models.Downvote)
}

func (s *CommentService) voteComment(ctx context.Context, commentID string, voteType models.VoteType) (*models.Comment, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, errors.New("unauthorized")
	}

	// Rate limiting
	if err := s.rateLimiter.Check(ctx, userID, "vote"); err != nil {
		return nil, err
	}

	// Get existing vote
	existingVote, err := s.voteRepo.GetCommentVote(ctx, userID, commentID)
	if err != nil {
		return nil, err
	}

	var upvoteDelta, downvoteDelta int

	if existingVote == nil {
		// New vote
		vote := &models.CommentVote{
			UserID:    userID,
			CommentID: commentID,
			VoteType:  voteType,
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		if err := s.voteRepo.CreateCommentVote(ctx, vote); err != nil {
			return nil, err
		}

		if voteType == models.Upvote {
			upvoteDelta = 1
		} else {
			downvoteDelta = 1
		}
	} else if existingVote.VoteType == voteType {
		// Remove vote
		if err := s.voteRepo.DeleteCommentVote(ctx, userID, commentID); err != nil {
			return nil, err
		}

		if voteType == models.Upvote {
			upvoteDelta = -1
		} else {
			downvoteDelta = -1
		}
	} else {
		// Change vote
		if err := s.voteRepo.UpdateCommentVote(ctx, userID, commentID, voteType); err != nil {
			return nil, err
		}

		if voteType == models.Upvote {
			upvoteDelta = 1
			downvoteDelta = -1
		} else {
			upvoteDelta = -1
			downvoteDelta = 1
		}
	}

	// Update comment vote counts
	if err := s.commentRepo.UpdateVoteCount(ctx, commentID, upvoteDelta, downvoteDelta); err != nil {
		return nil, err
	}

	// Return updated comment
	return s.commentRepo.GetByID(ctx, commentID)
}
