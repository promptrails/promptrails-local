package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"github.com/segmentio/ksuid"
)

type ScoreHandler struct {
	store *store.Store
}

func (h *ScoreHandler) List(c echo.Context) error {
	p := getPagination(c)
	scores, total := h.store.ListScores(p.Page, p.Limit)
	return listResponse(c, scores, total, p)
}

func (h *ScoreHandler) Create(c echo.Context) error {
	wid := getWorkspaceID()
	var req model.CreateScoreRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Name == "" || req.TraceID == "" {
		return badRequest(c, "name and trace_id are required")
	}

	now := time.Now()
	score := model.Score{
		ID:          ksuid.New().String(),
		WorkspaceID: wid,
		TraceID:     req.TraceID,
		SpanID:      req.SpanID,
		Name:        req.Name,
		Value:       req.Value,
		StringValue: req.StringValue,
		BoolValue:   req.BoolValue,
		DataType:    req.DataType,
		Comment:     req.Comment,
		Source:      req.Source,
		ConfigID:    req.ConfigID,
		ExecutionID: req.ExecutionID,
		AgentID:     req.AgentID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.store.CreateScore(score)
	return dataResponse(c, http.StatusCreated, score)
}

func (h *ScoreHandler) Get(c echo.Context) error {
	score, ok := h.store.GetScore(c.Param("scoreId"))
	if !ok {
		return notFound(c, "score not found")
	}
	return dataResponse(c, http.StatusOK, score)
}

func (h *ScoreHandler) Update(c echo.Context) error {
	score, ok := h.store.GetScore(c.Param("scoreId"))
	if !ok {
		return notFound(c, "score not found")
	}

	var req model.UpdateScoreRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if req.Value != nil {
		score.Value = req.Value
	}
	if req.StringValue != nil {
		score.StringValue = req.StringValue
	}
	if req.BoolValue != nil {
		score.BoolValue = req.BoolValue
	}
	if req.Comment != nil {
		score.Comment = req.Comment
	}
	score.UpdatedAt = time.Now()

	h.store.UpdateScore(score)
	return dataResponse(c, http.StatusOK, score)
}

func (h *ScoreHandler) Delete(c echo.Context) error {
	if !h.store.DeleteScore(c.Param("scoreId")) {
		return notFound(c, "score not found")
	}
	return c.NoContent(http.StatusNoContent)
}

// Score Configs

func (h *ScoreHandler) ListConfigs(c echo.Context) error {
	p := getPagination(c)
	configs, total := h.store.ListScoreConfigs(p.Page, p.Limit)
	return listResponse(c, configs, total, p)
}

func (h *ScoreHandler) CreateConfig(c echo.Context) error {
	wid := getWorkspaceID()
	var req model.CreateScoreConfigRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return badRequest(c, "name is required")
	}

	now := time.Now()
	cfg := model.ScoreConfig{
		ID:          ksuid.New().String(),
		WorkspaceID: wid,
		Name:        req.Name,
		DataType:    req.DataType,
		MinValue:    req.MinValue,
		MaxValue:    req.MaxValue,
		Categories:  req.Categories,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.store.CreateScoreConfig(cfg)
	return dataResponse(c, http.StatusCreated, cfg)
}

func (h *ScoreHandler) GetConfig(c echo.Context) error {
	cfg, ok := h.store.GetScoreConfig(c.Param("configId"))
	if !ok {
		return notFound(c, "score config not found")
	}
	return dataResponse(c, http.StatusOK, cfg)
}

func (h *ScoreHandler) UpdateConfig(c echo.Context) error {
	cfg, ok := h.store.GetScoreConfig(c.Param("configId"))
	if !ok {
		return notFound(c, "score config not found")
	}

	var req model.UpdateScoreConfigRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if req.Name != nil {
		cfg.Name = *req.Name
	}
	if req.Description != nil {
		cfg.Description = req.Description
	}
	if req.IsActive != nil {
		cfg.IsActive = *req.IsActive
	}
	if req.MinValue != nil {
		cfg.MinValue = req.MinValue
	}
	if req.MaxValue != nil {
		cfg.MaxValue = req.MaxValue
	}
	if req.Categories != nil {
		cfg.Categories = req.Categories
	}
	cfg.UpdatedAt = time.Now()

	h.store.UpdateScoreConfig(cfg)
	return dataResponse(c, http.StatusOK, cfg)
}

func (h *ScoreHandler) DeleteConfig(c echo.Context) error {
	if !h.store.DeleteScoreConfig(c.Param("configId")) {
		return notFound(c, "score config not found")
	}
	return c.NoContent(http.StatusNoContent)
}
