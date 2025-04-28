package services

import (
	"context"
	"errors"
	"skyphin-api/internal/models"
	"skyphin-api/internal/repositories"
	"skyphin-api/internal/utils"
	"time"
)

type VoteService struct {
	voteRepo    *repositories.VoteRepository
	urlService  *utils.URLService
	rateLimiter *utils.RateLimiter
}

func NewVoteService(
	voteRepo *repositories.VoteRepository,
	urlService *utils.URLService,
	rateLimiter *utils.RateLimiter,
) *VoteService {
	return &VoteService{
		voteRepo:    voteRepo,
		urlService:  urlService,
		rateLimiter: rateLimiter,
	}
}

func (s *VoteService) GetURLVoteCount(ctx context.Context, url string) (*models.URLVoteResult, error) {
	normalizedURL, err := s.urlService.Normalize(url)
	if err != nil {
		return nil, errors.New("invalid URL")
	}

	upvotes, downvotes, err := s.voteRepo.GetURLVoteCounts(ctx, normalizedURL)
	if err != nil {
		return nil, err
	}

	return &models.URLVoteResult{
		URL:           normalizedURL,
		UpvoteCount:   upvotes,
		DownvoteCount: downvotes,
	}, nil
}

func (s *VoteService) UpvoteURL(ctx context.Context, url string) (*models.URLVoteResult, error) {
	return s.voteURL(ctx, url, models.Upvote)
}

func (s *VoteService) DownvoteURL(ctx context.Context, url string) (*models.URLVoteResult, error) {
	return s.voteURL(ctx, url, models.Downvote)
}

func (s *VoteService) voteURL(ctx context.Context, url string, voteType models.VoteType) (*models.URLVoteResult, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return nil, errors.New("unauthorized")
	}

	// Rate limiting
	if err := s.rateLimiter.Check(ctx, userID, "vote"); err != nil {
		return nil, err
	}

	normalizedURL, err := s.urlService.Normalize(url)
	if err != nil {
		return nil, errors.New("invalid URL")
	}

	// Get existing vote
	existingVote, err := s.voteRepo.GetURLVote(ctx, userID, normalizedURL)
	if err != nil {
		return nil, err
	}

	if existingVote == nil {
		// New vote
		vote := &models.URLVote{
			UserID:    userID,
			URL:       normalizedURL,
			VoteType:  voteType,
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		if err := s.voteRepo.CreateURLVote(ctx, vote); err != nil {
			return nil, err
		}
	} else if existingVote.VoteType == voteType {
		// Remove vote
		if err := s.voteRepo.DeleteURLVote(ctx, userID, normalizedURL); err != nil {
			return nil, err
		}
	} else {
		// Change vote
		if err := s.voteRepo.UpdateURLVote(ctx, userID, normalizedURL, voteType); err != nil {
			return nil, err
		}
	}

	upvotes, downvotes, err := s.voteRepo.GetURLVoteCounts(ctx, normalizedURL)
	if err != nil {
		return nil, err
	}

	return &models.URLVoteResult{
		URL:           normalizedURL,
		UpvoteCount:   upvotes,
		DownvoteCount: downvotes,
	}, nil
}
