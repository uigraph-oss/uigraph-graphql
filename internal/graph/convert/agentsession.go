package convert

import (
	"encoding/json"

	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func rawToJSON(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}
	var out any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil
	}
	return out
}

func AgentSessionTotalsToModel(t uigraphapi.AgentSessionTotals) *model.AgentSessionTotals {
	return &model.AgentSessionTotals{
		StepCount:          t.StepCount,
		InputTokens:        t.InputTokens,
		OutputTokens:       t.OutputTokens,
		ReasoningTokens:    t.ReasoningTokens,
		CachedInputTokens:  t.CachedInputTokens,
		CachedOutputTokens: t.CachedOutputTokens,
		CostUsd:            t.CostUSD,
		UnpricedSteps:      t.UnpricedSteps,
		StepDurationMs:     t.StepDurationMs,
	}
}

func AgentSessionToModel(s uigraphapi.AgentSession, actors map[string]*uigraphapi.Actor) *model.AgentSession {
	out := &model.AgentSession{
		ID:               s.ID,
		OrgID:            s.OrgID,
		Type:             s.Type,
		Status:           s.Status,
		UserID:           s.UserID,
		ServiceAccountID: s.ServiceAccountID,
		Title:            s.Title,
		ModelName:        s.ModelName,
		Metadata:         rawToJSON(s.Metadata),
		Report:           s.Report,
		Error:            s.Error,
		StartedAt:        s.StartedAt,
		UpdatedAt:        s.UpdatedAt,
		CompletedAt:      s.CompletedAt,
		Totals:           AgentSessionTotalsToModel(s.Totals),
	}

	if s.CompletedAt != nil {
		ms := int(s.CompletedAt.Sub(s.StartedAt).Milliseconds())
		out.DurationMs = &ms
	}

	actorID := ""
	if s.UserID != nil {
		actorID = *s.UserID
	} else if s.ServiceAccountID != nil {
		actorID = *s.ServiceAccountID
	}
	if a := actors[actorID]; a != nil {
		name := a.Name
		out.ActorName = &name
		if a.AvatarURL != "" {
			avatar := a.AvatarURL
			out.ActorAvatarURL = &avatar
		}
	}
	return out
}

func AgentSessionListToModel(rows []uigraphapi.AgentSession, actors map[string]*uigraphapi.Actor) []*model.AgentSession {
	out := make([]*model.AgentSession, len(rows))
	for i, row := range rows {
		out[i] = AgentSessionToModel(row, actors)
	}
	return out
}

func AgentSessionStepToModel(st uigraphapi.AgentSessionStep) *model.AgentSessionStep {
	return &model.AgentSessionStep{
		ID:                 st.ID,
		SessionID:          st.SessionID,
		Seq:                st.Seq,
		Kind:               st.Kind,
		Name:               st.Name,
		ModelName:          st.ModelName,
		Input:              rawToJSON(st.Input),
		Output:             rawToJSON(st.Output),
		Text:               st.Text,
		FinishReason:       st.FinishReason,
		Error:              st.Error,
		InputTokens:        st.InputTokens,
		OutputTokens:       st.OutputTokens,
		ReasoningTokens:    st.ReasoningTokens,
		CachedInputTokens:  st.CachedInputTokens,
		CachedOutputTokens: st.CachedOutputTokens,
		CostUsd:            st.CostUSD,
		StartedAt:          st.StartedAt,
		CompletedAt:        st.CompletedAt,
		DurationMs:         int(st.CompletedAt.Sub(st.StartedAt).Milliseconds()),
	}
}

func AgentSessionStepListToModel(rows []uigraphapi.AgentSessionStep) []*model.AgentSessionStep {
	out := make([]*model.AgentSessionStep, len(rows))
	for i, row := range rows {
		out[i] = AgentSessionStepToModel(row)
	}
	return out
}

func AgentSessionSummaryToModel(s uigraphapi.AgentSessionSummary) *model.AgentSessionSummary {
	byType := make([]*model.AgentSessionTypeSummary, len(s.ByType))
	for i, t := range s.ByType {
		byType[i] = &model.AgentSessionTypeSummary{
			Type:              t.Type,
			TotalSessions:     t.TotalSessions,
			CompletedSessions: t.CompletedSessions,
			FailedSessions:    t.FailedSessions,
			RunningSessions:   t.RunningSessions,
			TotalDurationMs:   t.TotalDurationMs,
			Totals:            AgentSessionTotalsToModel(t.Totals),
		}
	}
	return &model.AgentSessionSummary{
		Period:            s.Period,
		TotalSessions:     s.TotalSessions,
		CompletedSessions: s.CompletedSessions,
		FailedSessions:    s.FailedSessions,
		RunningSessions:   s.RunningSessions,
		TotalDurationMs:   s.TotalDurationMs,
		Totals:            AgentSessionTotalsToModel(s.Totals),
		ByType:            byType,
	}
}
