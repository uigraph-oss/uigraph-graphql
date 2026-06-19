package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func FlowComponentFieldToModel(f uigraphapi.FlowDiagramComponentField) *model.FlowDiagramComponentField {
	return &model.FlowDiagramComponentField{
		FlowDiagramComponentFieldID: f.FlowDiagramComponentFieldID,
		Label:                       f.Label,
		Type:                        f.Type,
		Required:                    f.Required,
		Readonly:                    f.Readonly,
		Options:                     f.Options,
		Order:                       f.Order,
	}
}

func FlowComponentToModel(c uigraphapi.FlowDiagramComponent) *model.FlowDiagramComponent {
	fields := make([]*model.FlowDiagramComponentField, len(c.FlowDiagramComponentFields))
	for i, f := range c.FlowDiagramComponentFields {
		fields[i] = FlowComponentFieldToModel(f)
	}
	return &model.FlowDiagramComponent{
		ComponentID: c.ComponentID, Type: c.Type, Name: c.Name,
		Description: c.Description, Category: c.Category, Tags: c.Tags,
		Slug: c.Slug, PreviewImageJpg: c.PreviewImageJpg, IsActive: c.IsActive,
		Order: c.Order, OrganizationID: c.OrganizationID,
		FlowDiagramComponentFields: fields,
	}
}

func FlowComponentsToModel(components []uigraphapi.FlowDiagramComponent) []*model.FlowDiagramComponent {
	out := make([]*model.FlowDiagramComponent, len(components))
	for i, c := range components {
		out[i] = FlowComponentToModel(c)
	}
	return out
}

func ComponentFieldToModel(f uigraphapi.ComponentField) *model.ComponentField {
	return &model.ComponentField{
		ComponentFieldID: f.ComponentFieldID,
		Label:            f.Label,
		Type:             f.Type,
		Required:         f.Required,
		Readonly:         f.Readonly,
		Options:          f.Options,
		Order:            f.Order,
	}
}

func ComponentToModel(c uigraphapi.Component) *model.Component {
	fields := make([]*model.ComponentField, len(c.ComponentFields))
	for i, f := range c.ComponentFields {
		fields[i] = ComponentFieldToModel(f)
	}
	return &model.Component{
		ComponentID: c.ComponentID, Type: c.Type, Name: c.Name,
		Description: c.Description, Category: c.Category, Tags: c.Tags,
		Slug: c.Slug, PreviewImageJpg: c.PreviewImageJpg, IsActive: c.IsActive,
		Order: c.Order, ComponentFields: fields,
	}
}

func ComponentsToModel(components []uigraphapi.Component) []*model.Component {
	out := make([]*model.Component, len(components))
	for i, c := range components {
		out[i] = ComponentToModel(c)
	}
	return out
}
