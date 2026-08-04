package convert

import (
	"time"

	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func TimelineEventToModel(e *uigraphapi.TimelineEvent) *model.TimelineEvent {
	if e == nil {
		return nil
	}
	touches := make([]*model.TimelineTouch, len(e.Touches))
	for i, t := range e.Touches {
		touches[i] = &model.TimelineTouch{ID: t.ID, Label: t.Label, Kind: t.Kind}
	}
	return &model.TimelineEvent{
		ID:                 e.ID,
		OrgID:              e.OrgID,
		ServiceID:          e.ServiceID,
		Type:               e.Type,
		Title:              e.Title,
		Summary:            e.Summary,
		EventDate:          e.EventDate,
		Version:            e.Version,
		AdrNumber:          e.ADRNumber,
		DecisionStatus:     e.DecisionStatus,
		SourceLabel:        e.SourceLabel,
		SourceURL:          e.SourceURL,
		IsAgentSummarized:  e.IsAgentSummarized,
		Origin:             e.Origin,
		Touches:            touches,
		AttachmentAssetID:  e.AttachmentAssetID,
		AttachmentFileName: e.AttachmentFileName,
		AttachmentFileType: e.AttachmentFileType,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}

func TimelineEventListToModel(rows []uigraphapi.TimelineEvent) []*model.TimelineEvent {
	out := make([]*model.TimelineEvent, len(rows))
	for i := range rows {
		out[i] = TimelineEventToModel(&rows[i])
	}
	return out
}

func CreateTimelineEventInputToBody(in model.CreateTimelineEventInput) map[string]interface{} {
	return timelineEventInputToBody(in.Type, in.Title, in.Summary, in.EventDate, in.Version, in.AdrNumber,
		in.DecisionStatus, in.SourceLabel, in.SourceURL, in.IsAgentSummarized, in.Touches,
		in.AttachmentAssetID, in.AttachmentFileName, in.AttachmentFileType)
}

func UpdateTimelineEventInputToBody(in model.UpdateTimelineEventInput) map[string]interface{} {
	return timelineEventInputToBody(in.Type, in.Title, in.Summary, in.EventDate, in.Version, in.AdrNumber,
		in.DecisionStatus, in.SourceLabel, in.SourceURL, in.IsAgentSummarized, in.Touches,
		in.AttachmentAssetID, in.AttachmentFileName, in.AttachmentFileType)
}

func timelineEventInputToBody(
	eventType, title, summary string,
	eventDate time.Time,
	version, adrNumber, decisionStatus, sourceLabel, sourceURL *string,
	isAgentSummarized *bool,
	touches []*model.TimelineTouchInput,
	attachmentAssetID, attachmentFileName, attachmentFileType *string,
) map[string]interface{} {
	body := map[string]interface{}{
		"type":      eventType,
		"title":     title,
		"summary":   summary,
		"eventDate": eventDate,
	}
	setIfPresent(body, "version", version)
	setIfPresent(body, "adrNumber", adrNumber)
	setIfPresent(body, "decisionStatus", decisionStatus)
	setIfPresent(body, "sourceLabel", sourceLabel)
	setIfPresent(body, "sourceUrl", sourceURL)
	if isAgentSummarized != nil {
		body["isAgentSummarized"] = *isAgentSummarized
	}
	touchBodies := make([]map[string]interface{}, len(touches))
	for i, t := range touches {
		touchBodies[i] = map[string]interface{}{"id": t.ID, "label": t.Label, "kind": t.Kind}
	}
	body["touches"] = touchBodies
	setIfPresent(body, "attachmentAssetId", attachmentAssetID)
	setIfPresent(body, "attachmentFileName", attachmentFileName)
	setIfPresent(body, "attachmentFileType", attachmentFileType)
	return body
}
