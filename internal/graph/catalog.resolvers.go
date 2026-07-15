package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/generated"
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func (r *aPIGroupVersionResolver) CreatedByActor(ctx context.Context, obj *model.APIGroupVersion) (*model.Actor, error) {
	return r.resolveActor(ctx, obj.OrgID, obj.CreatedBy)
}

func (r *mutationResolver) CreateService(ctx context.Context, orgID string, input model.CreateServiceInput) (*model.Service, error) {
	s, err := r.Catalog.CreateService(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ServiceToModel(s), nil
}

func (r *mutationResolver) UpdateService(ctx context.Context, orgID string, id string, input model.UpdateServiceInput) (*model.Service, error) {
	s, err := r.Catalog.UpdateService(ctx, orgID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ServiceToModel(s), nil
}

func (r *mutationResolver) DeleteService(ctx context.Context, orgID string, id string) (bool, error) {
	return true, r.Catalog.DeleteService(ctx, orgID, id)
}

func (r *mutationResolver) CreateAPIGroup(ctx context.Context, orgID string, serviceID string, input model.CreateAPIGroupInput) (*model.APIGroup, error) {
	g, err := r.Catalog.CreateAPIGroup(ctx, orgID, serviceID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.APIGroupToModel(g), nil
}

func (r *mutationResolver) UpdateAPIGroup(ctx context.Context, orgID string, serviceID string, id string, input model.UpdateAPIGroupInput) (*model.APIGroup, error) {
	g, err := r.Catalog.UpdateAPIGroup(ctx, orgID, serviceID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.APIGroupToModel(g), nil
}

func (r *mutationResolver) DeleteAPIGroup(ctx context.Context, orgID string, serviceID string, id string) (bool, error) {
	return true, r.Catalog.DeleteAPIGroup(ctx, orgID, serviceID, id)
}

func (r *mutationResolver) SyncAPIGroup(ctx context.Context, orgID string, serviceID string, input model.SyncAPIGroupInput) (*model.SyncAPIGroupResult, error) {
	out, err := r.Catalog.SyncAPIGroup(ctx, orgID, serviceID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return &model.SyncAPIGroupResult{
		APIGroupID:     convert.StrFromMap(out, "apiGroupId"),
		VersionCreated: convert.BoolFromMap(out, "versionCreated"),
	}, nil
}

func (r *mutationResolver) RestoreAPIGroupVersion(ctx context.Context, orgID string, serviceID string, apiGroupID string, versionID string) (*model.APIGroup, error) {
	g, err := r.Catalog.RestoreAPIGroupVersion(ctx, orgID, serviceID, apiGroupID, versionID)
	if err != nil {
		return nil, err
	}
	return convert.APIGroupToModel(g), nil
}

func (r *mutationResolver) CreateServiceDoc(ctx context.Context, orgID string, serviceID string, input model.CreateServiceDocInput) (*model.ServiceDoc, error) {
	d, err := r.Catalog.CreateServiceDoc(ctx, orgID, serviceID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ServiceDocToModel(d), nil
}

func (r *mutationResolver) DeleteServiceDoc(ctx context.Context, orgID string, serviceID string, docID string) (bool, error) {
	return true, r.Catalog.DeleteServiceDoc(ctx, orgID, serviceID, docID)
}

func (r *mutationResolver) CreateServiceDiagram(ctx context.Context, orgID string, serviceID string, input model.CreateServiceDiagramInput) (*model.ServiceDiagram, error) {
	d, err := r.Catalog.CreateServiceDiagram(ctx, orgID, serviceID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ServiceDiagramToModel(d), nil
}

func (r *mutationResolver) DeleteServiceDiagram(ctx context.Context, orgID string, serviceID string, diagramID string) (bool, error) {
	return true, r.Catalog.DeleteServiceDiagram(ctx, orgID, serviceID, diagramID)
}

func (r *mutationResolver) CreateServiceDb(ctx context.Context, orgID string, serviceID string, input model.CreateServiceDBInput) (*model.ServiceDb, error) {
	d, err := r.Catalog.CreateServiceDB(ctx, orgID, serviceID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ServiceDBToModel(d), nil
}

func (r *mutationResolver) UpdateServiceDb(ctx context.Context, orgID string, serviceID string, id string, input model.UpdateServiceDBInput) (*model.ServiceDb, error) {
	d, err := r.Catalog.UpdateServiceDB(ctx, orgID, serviceID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ServiceDBToModel(d), nil
}

func (r *mutationResolver) DeleteServiceDb(ctx context.Context, orgID string, serviceID string, id string) (bool, error) {
	return true, r.Catalog.DeleteServiceDB(ctx, orgID, serviceID, id)
}

func (r *mutationResolver) CreateServiceDBVersion(ctx context.Context, orgID string, serviceID string, serviceDbID string, input model.CreateServiceDBVersionInput) (*model.ServiceDBVersion, error) {
	v, err := r.Catalog.CreateServiceDBVersion(ctx, orgID, serviceID, serviceDbID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ServiceDBVersionToModel(orgID, *v), nil
}

func (r *mutationResolver) RestoreServiceDBVersion(ctx context.Context, orgID string, serviceID string, serviceDbID string, versionID string) (*model.ServiceDb, error) {
	d, err := r.Catalog.RestoreServiceDBVersion(ctx, orgID, serviceID, serviceDbID, versionID)
	if err != nil {
		return nil, err
	}
	return convert.ServiceDBToModel(d), nil
}

func (r *mutationResolver) CreateAPIEndpoint(ctx context.Context, orgID string, serviceID string, apiGroupID string, input model.CreateAPIEndpointInput) (*model.APIEndpoint, error) {
	e, err := r.Catalog.CreateAPIEndpoint(ctx, orgID, serviceID, apiGroupID, convert.APIEndpointInputMap(input))
	if err != nil {
		return nil, err
	}
	return convert.APIEndpointToModel(e), nil
}

func (r *mutationResolver) UpdateAPIEndpoint(ctx context.Context, orgID string, serviceID string, apiGroupID string, id string, input model.UpdateAPIEndpointInput) (*model.APIEndpoint, error) {
	e, err := r.Catalog.UpdateAPIEndpoint(ctx, orgID, serviceID, apiGroupID, id, convert.APIEndpointInputMap(input))
	if err != nil {
		return nil, err
	}
	return convert.APIEndpointToModel(e), nil
}

func (r *mutationResolver) DeleteAPIEndpoint(ctx context.Context, orgID string, serviceID string, apiGroupID string, id string) (bool, error) {
	return true, r.Catalog.DeleteAPIEndpoint(ctx, orgID, serviceID, apiGroupID, id)
}

func (r *queryResolver) Services(ctx context.Context, orgID string, folderID *string, teamID *string, search *string, sortBy *string, sortDir *string, limit *int, offset *int) (*model.ServicePage, error) {
	p := uigraphapi.ListParams{
		FolderID: derefStr(folderID),
		TeamID:   derefStr(teamID),
		Search:   derefStr(search),
		SortBy:   derefStr(sortBy),
		SortDir:  derefStr(sortDir),
		Limit:    limit,
		Offset:   offset,
	}
	services, total, err := r.Catalog.ListServices(ctx, orgID, p)
	if err != nil {
		return nil, err
	}
	items := convert.ServicesToModel(services)

	stats, err := r.Catalog.ListServiceStats(ctx, orgID, nil)
	if err != nil {
		return nil, err
	}
	statsByID := make(map[string]*model.ServiceStats, len(stats))
	for _, s := range stats {
		statsByID[s.ServiceID] = convert.ServiceStatsToModel(s)
	}
	for _, item := range items {
		item.Stats = statsByID[item.ID]
	}

	return &model.ServicePage{Items: items, TotalCount: total}, nil
}

func (r *queryResolver) Service(ctx context.Context, orgID string, id string) (*model.Service, error) {
	s, err := r.Catalog.GetService(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	svc := convert.ServiceToModel(s)

	stats, err := r.Catalog.ListServiceStats(ctx, orgID, &id)
	if err != nil {
		return nil, err
	}
	if len(stats) > 0 {
		svc.Stats = convert.ServiceStatsToModel(stats[0])
	}

	return svc, nil
}

func (r *queryResolver) APIGroups(ctx context.Context, orgID string, serviceID string) ([]*model.APIGroup, error) {
	groups, err := r.Catalog.ListAPIGroups(ctx, orgID, serviceID)
	if err != nil {
		return nil, err
	}
	return convert.APIGroupsToModel(groups), nil
}

func (r *queryResolver) APIGroup(ctx context.Context, orgID string, serviceID string, id string) (*model.APIGroup, error) {
	g, err := r.Catalog.GetAPIGroup(ctx, orgID, serviceID, id)
	if err != nil {
		return nil, err
	}
	return convert.APIGroupToModel(g), nil
}

func (r *queryResolver) APIGroupVersions(ctx context.Context, orgID string, serviceID string, apiGroupID string) ([]*model.APIGroupVersion, error) {
	versions, err := r.Catalog.ListAPIGroupVersions(ctx, orgID, serviceID, apiGroupID)
	if err != nil {
		return nil, err
	}
	return convert.APIGroupVersionsToModel(orgID, versions), nil
}

func (r *queryResolver) ServiceDocs(ctx context.Context, orgID string, serviceID string) ([]*model.ServiceDoc, error) {
	docs, err := r.Catalog.ListServiceDocs(ctx, orgID, serviceID)
	if err != nil {
		return nil, err
	}
	return convert.ServiceDocsToModel(docs), nil
}

func (r *queryResolver) ServiceDiagrams(ctx context.Context, orgID string, serviceID string) ([]*model.ServiceDiagram, error) {
	diagrams, err := r.Catalog.ListServiceDiagrams(ctx, orgID, serviceID)
	if err != nil {
		return nil, err
	}
	return convert.ServiceDiagramsToModel(diagrams), nil
}

func (r *queryResolver) ServiceDBs(ctx context.Context, orgID string, serviceID string) ([]*model.ServiceDb, error) {
	dbs, err := r.Catalog.ListServiceDBs(ctx, orgID, serviceID)
	if err != nil {
		return nil, err
	}
	return convert.ServiceDBsToModel(dbs), nil
}

func (r *queryResolver) ServiceDb(ctx context.Context, orgID string, serviceID string, id string) (*model.ServiceDb, error) {
	d, err := r.Catalog.GetServiceDB(ctx, orgID, serviceID, id)
	if err != nil {
		return nil, err
	}
	return convert.ServiceDBToModel(d), nil
}

func (r *queryResolver) ServiceDBVersions(ctx context.Context, orgID string, serviceID string, serviceDbID string) ([]*model.ServiceDBVersion, error) {
	versions, err := r.Catalog.ListServiceDBVersions(ctx, orgID, serviceID, serviceDbID)
	if err != nil {
		return nil, err
	}
	return convert.ServiceDBVersionsToModel(orgID, versions), nil
}

func (r *queryResolver) APIEndpoints(ctx context.Context, orgID string, serviceID string, apiGroupID string, versionID *string) ([]*model.APIEndpoint, error) {
	vid := ""
	if versionID != nil {
		vid = *versionID
	}
	endpoints, err := r.Catalog.ListAPIEndpoints(ctx, orgID, serviceID, apiGroupID, vid)
	if err != nil {
		return nil, err
	}
	return convert.APIEndpointsToModel(endpoints), nil
}

func (r *queryResolver) APIEndpoint(ctx context.Context, orgID string, serviceID string, apiGroupID string, id string) (*model.APIEndpoint, error) {
	e, err := r.Catalog.GetAPIEndpoint(ctx, orgID, serviceID, apiGroupID, id)
	if err != nil {
		return nil, err
	}
	return convert.APIEndpointToModel(e), nil
}

func (r *queryResolver) APIEndpointByID(ctx context.Context, orgID string, id string) (*model.APIEndpoint, error) {
	e, err := r.Catalog.GetAPIEndpointByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return convert.APIEndpointToModel(e), nil
}

func (r *queryResolver) ServiceDocByID(ctx context.Context, orgID string, id string) (*model.ServiceDoc, error) {
	d, err := r.Catalog.GetServiceDocByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return convert.ServiceDocToModel(d), nil
}

func (r *queryResolver) APIGroupSpec(ctx context.Context, orgID string, serviceID string, apiGroupID string, versionID *string) (*model.FileDownload, error) {
	vid := ""
	if versionID != nil {
		vid = *versionID
	}
	spec, err := r.Catalog.GetAPIGroupSpec(ctx, orgID, serviceID, apiGroupID, vid)
	if err != nil {
		return nil, err
	}
	return &model.FileDownload{
		APIGroupID: spec.APIGroupID,
		FileName:   spec.FileName,
		Content:    spec.Content,
	}, nil
}

func (r *serviceResolver) CreatedByActor(ctx context.Context, obj *model.Service) (*model.Actor, error) {
	return r.resolveActor(ctx, obj.OrgID, obj.CreatedBy)
}

func (r *serviceResolver) UpdatedByActor(ctx context.Context, obj *model.Service) (*model.Actor, error) {
	if obj.UpdatedBy == nil {
		return nil, nil
	}
	return r.resolveActor(ctx, obj.OrgID, *obj.UpdatedBy)
}

func (r *serviceDBResolver) CreatedByActor(ctx context.Context, obj *model.ServiceDb) (*model.Actor, error) {
	return r.resolveActor(ctx, obj.OrgID, obj.CreatedBy)
}

func (r *serviceDBResolver) UpdatedByActor(ctx context.Context, obj *model.ServiceDb) (*model.Actor, error) {
	if obj.UpdatedBy == nil {
		return nil, nil
	}
	return r.resolveActor(ctx, obj.OrgID, *obj.UpdatedBy)
}

func (r *serviceDBVersionResolver) CreatedByActor(ctx context.Context, obj *model.ServiceDBVersion) (*model.Actor, error) {
	return r.resolveActor(ctx, obj.OrgID, obj.CreatedBy)
}

func (r *Resolver) APIGroupVersion() generated.APIGroupVersionResolver {
	return &aPIGroupVersionResolver{r}
}

func (r *Resolver) Service() generated.ServiceResolver { return &serviceResolver{r} }

func (r *Resolver) ServiceDB() generated.ServiceDBResolver { return &serviceDBResolver{r} }

func (r *Resolver) ServiceDBVersion() generated.ServiceDBVersionResolver {
	return &serviceDBVersionResolver{r}
}

type aPIGroupVersionResolver struct{ *Resolver }
type serviceResolver struct{ *Resolver }
type serviceDBResolver struct{ *Resolver }
type serviceDBVersionResolver struct{ *Resolver }
