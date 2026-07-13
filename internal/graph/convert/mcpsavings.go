package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func SavingsSummaryToModel(s *uigraphapi.SavingsSummary) *model.SavingsSummary {
	if s == nil {
		return nil
	}
	var modelID *string
	if s.ModelID != "" {
		modelID = &s.ModelID
	}
	return &model.SavingsSummary{
		OrgID:             s.OrgID,
		Period:            s.Period,
		ModelID:           modelID,
		TotalCalls:        s.TotalCalls,
		TotalTokensServed: s.TotalTokensServed,
		TotalTokensSaved:  s.TotalTokensSaved,
		CostServedUsd:     s.CostServedUSD,
		CostRawUsd:        s.CostRawUSD,
		CostSavedUsd:      s.CostSavedUSD,
		UniqueUsersCount:  s.UniqueUsersCount,
		TotalDurationMs:   s.TotalDurationMs,
		EstAgentTimeMs:    s.EstAgentTimeMs,
		TimeSavedMs:       s.TimeSavedMs,
	}
}

func DailySavingsToModel(d uigraphapi.DailySavings) *model.DailySavings {
	return &model.DailySavings{
		Date:              d.Date,
		TotalCalls:        d.TotalCalls,
		TotalTokensServed: d.TotalTokensServed,
		TotalTokensSaved:  d.TotalTokensSaved,
		CostServedUsd:     d.CostServedUSD,
		CostRawUsd:        d.CostRawUSD,
		CostSavedUsd:      d.CostSavedUSD,
		TotalDurationMs:   d.TotalDurationMs,
		EstAgentTimeMs:    d.EstAgentTimeMs,
		TimeSavedMs:       d.TimeSavedMs,
	}
}

func DailySavingsListToModel(rows []uigraphapi.DailySavings) []*model.DailySavings {
	out := make([]*model.DailySavings, len(rows))
	for i, row := range rows {
		out[i] = DailySavingsToModel(row)
	}
	return out
}

func ToolSavingsToModel(s uigraphapi.ToolSavings) *model.ToolSavings {
	return &model.ToolSavings{
		ToolName:        s.ToolName,
		TotalCalls:      s.TotalCalls,
		TokensSaved:     s.TokensSaved,
		CostSavedUsd:    s.CostSavedUSD,
		TotalDurationMs: s.TotalDurationMs,
		EstAgentTimeMs:  s.EstAgentTimeMs,
		TimeSavedMs:     s.TimeSavedMs,
	}
}

func ToolSavingsListToModel(rows []uigraphapi.ToolSavings) []*model.ToolSavings {
	out := make([]*model.ToolSavings, len(rows))
	for i, row := range rows {
		out[i] = ToolSavingsToModel(row)
	}
	return out
}

func ClientSavingsToModel(s uigraphapi.ClientSavings) *model.ClientSavings {
	return &model.ClientSavings{
		ClientName:      s.ClientName,
		TotalCalls:      s.TotalCalls,
		TokensSaved:     s.TokensSaved,
		CostSavedUsd:    s.CostSavedUSD,
		TotalDurationMs: s.TotalDurationMs,
	}
}

func ClientSavingsListToModel(rows []uigraphapi.ClientSavings) []*model.ClientSavings {
	out := make([]*model.ClientSavings, len(rows))
	for i, row := range rows {
		out[i] = ClientSavingsToModel(row)
	}
	return out
}

func ModelSavingsToModel(s uigraphapi.ModelSavings) *model.ModelSavings {
	return &model.ModelSavings{
		ModelID:      s.ModelID,
		DisplayName:  s.DisplayName,
		Provider:     s.Provider,
		TotalCalls:   s.TotalCalls,
		TokensSaved:  s.TokensSaved,
		CostRawUsd:   s.CostRawUSD,
		CostSavedUsd: s.CostSavedUSD,
	}
}

func ModelSavingsListToModel(rows []uigraphapi.ModelSavings) []*model.ModelSavings {
	out := make([]*model.ModelSavings, len(rows))
	for i, row := range rows {
		out[i] = ModelSavingsToModel(row)
	}
	return out
}

// UserSavingsToModel resolves DisplayName from actors (keyed by user ID or
// service account ID), falling back to "Service Account" or "Unknown User"
// when no actor was resolved for that ID.
func UserSavingsToModel(s uigraphapi.UserSavings, actors map[string]*uigraphapi.Actor) *model.UserSavings {
	id := ""
	if s.UserID != nil {
		id = *s.UserID
	} else if s.ServiceAccountID != nil {
		id = *s.ServiceAccountID
	}
	displayName := "Unknown User"
	var avatarURL *string
	if a := actors[id]; a != nil {
		displayName = a.Name
		if a.AvatarURL != "" {
			avatarURL = &a.AvatarURL
		}
	} else if s.ServiceAccountID != nil {
		displayName = "Service Account"
	}
	return &model.UserSavings{
		UserID:           s.UserID,
		ServiceAccountID: s.ServiceAccountID,
		DisplayName:      displayName,
		AvatarURL:        avatarURL,
		TotalCalls:       s.TotalCalls,
		TokensSaved:      s.TokensSaved,
		CostSavedUsd:     s.CostSavedUSD,
		TotalDurationMs:  s.TotalDurationMs,
	}
}

func UserSavingsListToModel(rows []uigraphapi.UserSavings, actors map[string]*uigraphapi.Actor) []*model.UserSavings {
	out := make([]*model.UserSavings, len(rows))
	for i, row := range rows {
		out[i] = UserSavingsToModel(row, actors)
	}
	return out
}
