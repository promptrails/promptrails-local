package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"github.com/segmentio/ksuid"
)

type MemoryHandler struct {
	store *store.Store
}

func (h *MemoryHandler) List(c echo.Context) error {
	agentID := c.Param("agentId")
	p := getPagination(c)
	memories, total := h.store.ListMemories(agentID, p.Page, p.Limit)
	return listResponse(c, memories, total, p)
}

func (h *MemoryHandler) Create(c echo.Context) error {
	agentID := c.Param("agentId")
	wid := getWorkspaceID()

	var req model.CreateMemoryRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Content == "" {
		return badRequest(c, "content is required")
	}

	now := time.Now()
	memory := model.AgentMemory{
		ID:          ksuid.New().String(),
		WorkspaceID: wid,
		AgentID:     agentID,
		Content:     req.Content,
		Metadata:    req.Metadata,
		MemoryType:  req.MemoryType,
		Importance:  req.Importance,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.store.CreateMemory(memory)
	return dataResponse(c, http.StatusCreated, memory)
}

func (h *MemoryHandler) Search(c echo.Context) error {
	agentID := c.Param("agentId")

	var req model.SearchMemoryRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	// Simple substring search across memories
	all, _ := h.store.ListMemories(agentID, 1, 1000)
	var results []model.AgentMemory
	for _, m := range all {
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(req.Query)) {
			results = append(results, m)
			if len(results) >= limit {
				break
			}
		}
	}
	if results == nil {
		results = []model.AgentMemory{}
	}
	return dataResponse(c, http.StatusOK, results)
}

func (h *MemoryHandler) Get(c echo.Context) error {
	memory, ok := h.store.GetMemory(c.Param("memoryId"))
	if !ok {
		return notFound(c, "memory not found")
	}
	return dataResponse(c, http.StatusOK, memory)
}

func (h *MemoryHandler) Delete(c echo.Context) error {
	if !h.store.DeleteMemory(c.Param("memoryId")) {
		return notFound(c, "memory not found")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *MemoryHandler) DeleteAll(c echo.Context) error {
	agentID := c.Param("agentId")
	count := h.store.DeleteAllMemories(agentID)
	return dataResponse(c, http.StatusOK, map[string]int{"deleted": count})
}
