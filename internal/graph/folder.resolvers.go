package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/model"
)

func (r *mutationResolver) CreateFolder(ctx context.Context, orgID string, input model.CreateFolderInput) (*model.Folder, error) {
	f, err := r.FolderAPI.CreateFolder(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.FolderToModel(f), nil
}

func (r *mutationResolver) UpdateFolder(ctx context.Context, orgID string, id string, input model.UpdateFolderInput) (*model.Folder, error) {
	f, err := r.FolderAPI.UpdateFolder(ctx, orgID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.FolderToModel(f), nil
}

func (r *mutationResolver) DeleteFolder(ctx context.Context, orgID string, id string) (bool, error) {
	return true, r.FolderAPI.DeleteFolder(ctx, orgID, id)
}

func (r *queryResolver) Folders(ctx context.Context, orgID string, typeArg *string, parentID *string) ([]*model.Folder, error) {
	t := ""
	if typeArg != nil {
		t = *typeArg
	}
	p := ""
	if parentID != nil {
		p = *parentID
	}
	folders, err := r.FolderAPI.ListFolders(ctx, orgID, t, p)
	if err != nil {
		return nil, err
	}
	return convert.FoldersToModel(folders), nil
}

func (r *queryResolver) Folder(ctx context.Context, orgID string, id string) (*model.Folder, error) {
	f, err := r.FolderAPI.GetFolder(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return convert.FolderToModel(f), nil
}
