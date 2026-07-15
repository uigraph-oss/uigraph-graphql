package convert

import (
	"github.com/uigraph/graphql/internal/graph/model"
	"github.com/uigraph/graphql/internal/uigraphapi"
)

func ChatSessionToModel(s *uigraphapi.ChatSession) *model.ChatSession {
	return &model.ChatSession{
		ID: s.ID, OrgID: s.OrgID, OwnerUserID: s.OwnerUserID, Title: s.Title,
		IsPinned: s.IsPinned, MessageCount: s.MessageCount,
		CreatedBy: s.CreatedBy, UpdatedBy: s.UpdatedBy,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

func ChatSessionsToModel(sessions []uigraphapi.ChatSession) []*model.ChatSession {
	out := make([]*model.ChatSession, len(sessions))
	for i := range sessions {
		out[i] = ChatSessionToModel(&sessions[i])
	}
	return out
}

func ChatMessageToModel(m *uigraphapi.ChatMessage) *model.ChatMessage {
	return &model.ChatMessage{
		ID: m.ID, OrgID: m.OrgID, ChatSessionID: m.ChatSessionID,
		Role: m.Role, Content: m.Content, Parts: m.Parts, CreatedAt: m.CreatedAt,
	}
}

func ChatMessagesToModel(messages []uigraphapi.ChatMessage) []*model.ChatMessage {
	out := make([]*model.ChatMessage, len(messages))
	for i := range messages {
		out[i] = ChatMessageToModel(&messages[i])
	}
	return out
}
