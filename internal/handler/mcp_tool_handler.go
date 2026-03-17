package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"github.com/segmentio/ksuid"
)

type MCPToolHandler struct {
	store *store.Store
}

func (h *MCPToolHandler) List(c echo.Context) error {
	p := getPagination(c)
	tools, total := h.store.ListMCPTools(p.Page, p.Limit)
	return listResponse(c, tools, total, p)
}

func (h *MCPToolHandler) Create(c echo.Context) error {
	wid := getWorkspaceID()
	var req model.CreateMCPToolRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return badRequest(c, "name is required")
	}

	now := time.Now()
	tool := model.MCPTool{
		ID:           ksuid.New().String(),
		WorkspaceID:  wid,
		Name:         req.Name,
		Description:  req.Description,
		Type:         req.Type,
		Config:       req.Config,
		Schema:       req.Schema,
		CredentialID: req.CredentialID,
		TemplateID:   req.TemplateID,
		Status:       "active",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	h.store.CreateMCPTool(tool)
	return dataResponse(c, http.StatusCreated, tool)
}

func (h *MCPToolHandler) Get(c echo.Context) error {
	tool, ok := h.store.GetMCPTool(c.Param("toolId"))
	if !ok {
		return notFound(c, "MCP tool not found")
	}
	return dataResponse(c, http.StatusOK, tool)
}

func (h *MCPToolHandler) Update(c echo.Context) error {
	tool, ok := h.store.GetMCPTool(c.Param("toolId"))
	if !ok {
		return notFound(c, "MCP tool not found")
	}

	var req model.UpdateMCPToolRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if req.Name != nil {
		tool.Name = *req.Name
	}
	if req.Description != nil {
		tool.Description = *req.Description
	}
	if req.Config != nil {
		tool.Config = req.Config
	}
	if req.Schema != nil {
		tool.Schema = req.Schema
	}
	if req.IsActive != nil {
		tool.IsActive = *req.IsActive
	}
	tool.UpdatedAt = time.Now()

	h.store.UpdateMCPTool(tool)
	return dataResponse(c, http.StatusOK, tool)
}

func (h *MCPToolHandler) Delete(c echo.Context) error {
	if !h.store.DeleteMCPTool(c.Param("toolId")) {
		return notFound(c, "MCP tool not found")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *MCPToolHandler) ListTemplates(c echo.Context) error {
	templates := h.store.ListMCPTemplates()
	return dataResponse(c, http.StatusOK, templates)
}

func (h *MCPToolHandler) GetTemplate(c echo.Context) error {
	tmpl, ok := h.store.GetMCPTemplate(c.Param("templateId"))
	if !ok {
		return notFound(c, "MCP template not found")
	}
	return dataResponse(c, http.StatusOK, tmpl)
}
