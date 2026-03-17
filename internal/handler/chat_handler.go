package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/fake"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"github.com/segmentio/ksuid"
)

type ChatHandler struct {
	store *store.Store
}

func (h *ChatHandler) ListSessions(c echo.Context) error {
	p := getPagination(c)
	sessions, total := h.store.ListChatSessions(p.Page, p.Limit)
	return listResponse(c, sessions, total, p)
}

func (h *ChatHandler) CreateSession(c echo.Context) error {
	wid := getWorkspaceID()
	var req model.CreateChatSessionRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.AgentID == "" {
		return badRequest(c, "agent_id is required")
	}

	now := time.Now()
	metadata := req.Metadata
	if metadata == nil {
		metadata = json.RawMessage("{}")
	}

	session := model.ChatSession{
		ID:          ksuid.New().String(),
		WorkspaceID: wid,
		AgentID:     req.AgentID,
		Title:       req.Title,
		Metadata:    metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.store.CreateChatSession(session)
	return dataResponse(c, http.StatusCreated, session)
}

func (h *ChatHandler) GetSession(c echo.Context) error {
	session, ok := h.store.GetChatSession(c.Param("sessionId"))
	if !ok {
		return notFound(c, "session not found")
	}
	return dataResponse(c, http.StatusOK, session)
}

func (h *ChatHandler) DeleteSession(c echo.Context) error {
	if !h.store.DeleteChatSession(c.Param("sessionId")) {
		return notFound(c, "session not found")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *ChatHandler) ListMessages(c echo.Context) error {
	sessionID := c.Param("sessionId")
	if _, ok := h.store.GetChatSession(sessionID); !ok {
		return notFound(c, "session not found")
	}
	p := getPagination(c)
	messages, total := h.store.ListChatMessages(sessionID, p.Page, p.Limit)
	return listResponse(c, messages, total, p)
}

func (h *ChatHandler) SendMessage(c echo.Context) error {
	sessionID := c.Param("sessionId")
	session, ok := h.store.GetChatSession(sessionID)
	if !ok {
		return notFound(c, "session not found")
	}

	var req model.SendMessageRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Content == "" {
		return badRequest(c, "content is required")
	}

	now := time.Now()

	// Create user message
	userMsg := model.ChatMessage{
		ID:        ksuid.New().String(),
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Content,
		Metadata:  json.RawMessage("{}"),
		CreatedAt: now,
	}
	h.store.CreateChatMessage(userMsg)

	// Generate fake assistant response
	assistantContent := fake.GenerateChatResponse(session.AgentID, req.Content)

	promptTokens := 100
	completionTokens := 60
	tokenCount := promptTokens + completionTokens
	cost := 0.0016

	assistantMsg := model.ChatMessage{
		ID:               ksuid.New().String(),
		SessionID:        sessionID,
		Role:             "assistant",
		Content:          assistantContent,
		Metadata:         json.RawMessage("{}"),
		Model:            "gpt-4o-mini",
		Cost:             &cost,
		TokenCount:       &tokenCount,
		PromptTokens:     &promptTokens,
		CompletionTokens: &completionTokens,
		CreatedAt:        now.Add(time.Millisecond * 200),
	}
	h.store.CreateChatMessage(assistantMsg)

	resp := model.SendMessageResponse{
		UserMessage:      &userMsg,
		AssistantMessage: &assistantMsg,
	}
	return dataResponse(c, http.StatusCreated, resp)
}
