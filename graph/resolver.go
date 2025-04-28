package graph

import (
	"context"
	"skyphin-api/internal/models"
	"skyphin-api/internal/services"
	"time"
)

type Resolver struct {
	commentService services.CommentService
	voteService    services.VoteService
	userService    services.UserService
}

func NewResolver(
	commentService services.CommentService,
	voteService services.VoteService,
	userService services.UserService,
) *Resolver {
	return &Resolver{
		commentService: commentService,
		voteService:    voteService,
		userService:    userService,
	}
}

// Root resolver methods
func (r *Resolver) Comment() CommentResolver {
	return &commentResolver{r}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

// Implement resolver types
type commentResolver struct{ *Resolver }

func (r *commentResolver) Author(ctx context.Context, obj *models.Comment) (*models.User, error) {
	return r.userService.GetUserByID(obj.AuthorID)
}

func (r *commentResolver) CreatedAt(ctx context.Context, obj *models.Comment) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

func (r *commentResolver) UpvoteCount(ctx context.Context, obj *models.Comment) (int, error) {
	return obj.UpvoteCount, nil
}

func (r *commentResolver) DownvoteCount(ctx context.Context, obj *models.Comment) (int, error) {
	return obj.DownvoteCount, nil
}

func (r *commentResolver) ReplyCount(ctx context.Context, obj *models.Comment) (int, error) {
	return obj.ReplyCount, nil
}

func (r *commentResolver) Replies(ctx context.Context, obj *models.Comment, first *int, after *string) (*models.CommentConnection, error) {
	limit := 10 // default
	if first != nil {
		limit = *first
	}
	return r.commentService.GetReplies(ctx, obj.ID, limit, after)
}

// type commentConnectionResolver struct{ *Resolver }

// func (r *commentConnectionResolver) TotalCount(ctx context.Context, obj *models.CommentConnection) (int, error) {
// 	return obj.TotalCount, nil
// }

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) AddComment(ctx context.Context, input models.AddCommentInput) (*models.Comment, error) {
	return r.commentService.CreateComment(ctx, &input)
}

func (r *mutationResolver) AddReply(ctx context.Context, input models.AddReplyInput) (*models.Comment, error) {
	return r.commentService.CreateReply(ctx, &input)
}

func (r *mutationResolver) UpvoteComment(ctx context.Context, commentID string) (*models.Comment, error) {
	return r.commentService.UpvoteComment(ctx, commentID)
}

func (r *mutationResolver) DownvoteComment(ctx context.Context, commentID string) (*models.Comment, error) {
	return r.commentService.DownvoteComment(ctx, commentID)
}

func (r *mutationResolver) UpvoteURL(ctx context.Context, url string) (*models.URLVoteResult, error) {
	return r.voteService.UpvoteURL(ctx, url)
}

func (r *mutationResolver) DownvoteURL(ctx context.Context, url string) (*models.URLVoteResult, error) {
	return r.voteService.DownvoteURL(ctx, url)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Comments(ctx context.Context, url string, first *int, after *string) (*models.CommentConnection, error) {
	limit := 10 // default
	if first != nil {
		limit = *first
	}
	return r.commentService.GetCommentsByURL(ctx, url, limit, after)
}

func (r *queryResolver) Comment(ctx context.Context, id string) (*models.Comment, error) {
	return r.commentService.GetComment(ctx, id)
}

func (r *queryResolver) Votes(ctx context.Context, url string) (*models.URLVoteResult, error) {
	return r.voteService.GetURLVoteCount(ctx, url)
}

func (r *queryResolver) CommentsCount(ctx context.Context, url string) (int, error) {
	return r.commentService.CountCommentsByURL(ctx, url)
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) CommentAdded(ctx context.Context, url string) (<-chan *models.Comment, error) {
	panic("implement with your pubsub system")
}

func (r *subscriptionResolver) TypingComment(ctx context.Context, url string) (<-chan *models.Comment, error) {
	panic("implement with your websocket service")
}

func (r *subscriptionResolver) ReplyAdded(ctx context.Context, url string, commentID string) (<-chan *models.Comment, error) {
	panic("implement with your pubsub system")
}

func (r *subscriptionResolver) TypingReply(ctx context.Context, url string, commentID string) (<-chan *models.Comment, error) {
	panic("implement with your websocket service")
}

func (r *subscriptionResolver) CommentVoted(ctx context.Context, commentID string) (<-chan *models.Comment, error) {
	panic("implement with your pubsub system")
}

func (r *subscriptionResolver) URLVoted(ctx context.Context, url string) (<-chan *models.URLVoteResult, error) {
	panic("implement with your pubsub system")
}
