package graph

import (
	"context"

	"github.com/uigraph/graphql/internal/graph/convert"
	"github.com/uigraph/graphql/internal/graph/model"
)

func (r *mutationResolver) CreateChatSession(ctx context.Context, orgID string, input model.CreateChatSessionInput) (*model.ChatSession, error) {
	s, err := r.Chat.CreateChatSession(ctx, orgID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ChatSessionToModel(s), nil
}

func (r *mutationResolver) UpdateChatSession(ctx context.Context, orgID string, id string, input model.UpdateChatSessionInput) (*model.ChatSession, error) {
	s, err := r.Chat.UpdateChatSession(ctx, orgID, id, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ChatSessionToModel(s), nil
}

func (r *mutationResolver) DeleteChatSession(ctx context.Context, orgID string, id string) (bool, error) {
	return true, r.Chat.DeleteChatSession(ctx, orgID, id)
}

func (r *mutationResolver) CreateChatMessage(ctx context.Context, orgID string, sessionID string, input model.CreateChatMessageInput) (*model.ChatMessage, error) {
	m, err := r.Chat.CreateChatMessage(ctx, orgID, sessionID, convert.ToMap(input))
	if err != nil {
		return nil, err
	}
	return convert.ChatMessageToModel(m), nil
}

func (r *queryResolver) ChatSessions(ctx context.Context, orgID string) ([]*model.ChatSession, error) {
	sessions, err := r.Chat.ListChatSessions(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return convert.ChatSessionsToModel(sessions), nil
}

func (r *queryResolver) ChatSession(ctx context.Context, orgID string, id string) (*model.ChatSessionWithMessages, error) {
	s, messages, err := r.Chat.GetChatSession(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	return &model.ChatSessionWithMessages{
		Session:  convert.ChatSessionToModel(s),
		Messages: convert.ChatMessagesToModel(messages),
	}, nil
}

func (r *queryResolver) ChatMessages(ctx context.Context, orgID string, sessionID string) ([]*model.ChatMessage, error) {
	messages, err := r.Chat.ListChatMessages(ctx, orgID, sessionID)
	if err != nil {
		return nil, err
	}
	return convert.ChatMessagesToModel(messages), nil
}
