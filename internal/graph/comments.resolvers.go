package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/graph/model"
)

func (r *commentResolver) CreatedByActor(ctx context.Context, obj *model.Comment) (*model.Actor, error) {
	return r.resolveActor(ctx, obj.OrgID, obj.CreatedBy)
}

func (r *mutationResolver) CreateComment(ctx context.Context, orgID string, input model.CreateCommentInput) (*model.Comment, error) {
	c, err := r.CommentAPI.CreateComment(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.CommentToModel(c), nil
}

func (r *mutationResolver) UpdateComment(ctx context.Context, orgID string, id string, input model.UpdateCommentInput) (*model.Comment, error) {
	c, err := r.CommentAPI.UpdateComment(ctx, orgID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.CommentToModel(c), nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, orgID string, id string) (bool, error) {
	return true, r.CommentAPI.DeleteComment(ctx, orgID, id)
}

func (r *queryResolver) Comments(ctx context.Context, orgID string, resourceID string) ([]*model.Comment, error) {
	cs, err := r.CommentAPI.ListComments(ctx, orgID, resourceID)
	if err != nil {
		return nil, err
	}
	return convert.CommentsToModel(cs), nil
}

func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

type commentResolver struct{ *Resolver }
