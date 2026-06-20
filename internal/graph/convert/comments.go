package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func CommentToModel(c *uigraphapi.Comment) *model.Comment {
	return &model.Comment{
		ID: c.ID, OrgID: c.OrgID, ResourceID: c.ResourceID,
		ParentCommentID: c.ParentCommentID, Text: c.Text,
		CreatedBy: c.CreatedBy, UpdatedBy: c.UpdatedBy,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}

func CommentsToModel(cs []uigraphapi.Comment) []*model.Comment {
	out := make([]*model.Comment, len(cs))
	for i := range cs {
		out[i] = CommentToModel(&cs[i])
	}
	return out
}
