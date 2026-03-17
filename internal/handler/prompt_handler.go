package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/fake"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
)

type PromptHandler struct {
	store  *store.Store
	logger *zap.Logger
}

func (h *PromptHandler) List(c echo.Context) error {
	p := getPagination(c)
	prompts, total := h.store.ListPrompts(p.Page, p.Limit)
	return listResponse(c, prompts, total, p)
}

func (h *PromptHandler) Create(c echo.Context) error {
	wid := getWorkspaceID()
	var req model.CreatePromptRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return badRequest(c, "name is required")
	}

	now := time.Now()
	prompt := model.Prompt{
		ID:          ksuid.New().String(),
		WorkspaceID: wid,
		Name:        req.Name,
		Description: req.Description,
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.store.CreatePrompt(prompt)
	return dataResponse(c, http.StatusCreated, prompt)
}

func (h *PromptHandler) Get(c echo.Context) error {
	prompt, ok := h.store.GetPrompt(c.Param("promptId"))
	if !ok {
		return notFound(c, "prompt not found")
	}
	return dataResponse(c, http.StatusOK, prompt)
}

func (h *PromptHandler) Update(c echo.Context) error {
	prompt, ok := h.store.GetPrompt(c.Param("promptId"))
	if !ok {
		return notFound(c, "prompt not found")
	}

	var req model.UpdatePromptRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if req.Name != nil {
		prompt.Name = *req.Name
	}
	if req.Description != nil {
		prompt.Description = *req.Description
	}
	if req.Status != nil {
		prompt.Status = *req.Status
	}
	prompt.UpdatedAt = time.Now()

	h.store.UpdatePrompt(prompt)
	return dataResponse(c, http.StatusOK, prompt)
}

func (h *PromptHandler) Delete(c echo.Context) error {
	if !h.store.DeletePrompt(c.Param("promptId")) {
		return notFound(c, "prompt not found")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PromptHandler) ListVersions(c echo.Context) error {
	promptID := c.Param("promptId")
	if _, ok := h.store.GetPrompt(promptID); !ok {
		return notFound(c, "prompt not found")
	}
	versions := h.store.ListPromptVersions(promptID)
	return dataResponse(c, http.StatusOK, versions)
}

func (h *PromptHandler) CreateVersion(c echo.Context) error {
	promptID := c.Param("promptId")
	if _, ok := h.store.GetPrompt(promptID); !ok {
		return notFound(c, "prompt not found")
	}

	var req model.CreatePromptVersionRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	version := model.PromptVersion{
		ID:                 ksuid.New().String(),
		PromptID:           promptID,
		Version:            req.Version,
		SystemPrompt:       req.SystemPrompt,
		UserPrompt:         req.UserPrompt,
		LLMModelID:         req.LLMModelID,
		FallbackLLMModelID: req.FallbackLLMModelID,
		Temperature:        req.Temperature,
		MaxTokens:          req.MaxTokens,
		TopP:               req.TopP,
		InputSchema:        req.InputSchema,
		OutputSchema:       req.OutputSchema,
		IsCurrent:          req.SetCurrent,
		Message:            req.Message,
		Config:             req.Config,
		CacheTimeout:       req.CacheTimeout,
		CreatedAt:          time.Now(),
	}

	if req.SetCurrent {
		h.store.DemotePromptVersions(promptID)
	}

	h.store.CreatePromptVersion(version)
	return dataResponse(c, http.StatusCreated, version)
}

func (h *PromptHandler) PromoteVersion(c echo.Context) error {
	promptID := c.Param("promptId")
	versionID := c.Param("versionId")

	version, ok := h.store.GetPromptVersion(versionID)
	if !ok || version.PromptID != promptID {
		return notFound(c, "version not found")
	}

	h.store.DemotePromptVersions(promptID)
	version.IsCurrent = true
	h.store.UpdatePromptVersion(version)
	return dataResponse(c, http.StatusOK, version)
}

func (h *PromptHandler) Preview(c echo.Context) error {
	promptID := c.Param("promptId")
	prompt, ok := h.store.GetPrompt(promptID)
	if !ok {
		return notFound(c, "prompt not found")
	}

	var req model.PreviewPromptRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	result := fake.GeneratePromptRunResponse(prompt.Name)
	return dataResponse(c, http.StatusOK, result)
}

func (h *PromptHandler) Run(c echo.Context) error {
	promptID := c.Param("promptId")
	prompt, ok := h.store.GetPrompt(promptID)
	if !ok {
		return notFound(c, "prompt not found")
	}

	var req model.RunPromptRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	result := fake.GeneratePromptRunResponse(prompt.Name)
	return dataResponse(c, http.StatusOK, result)
}
