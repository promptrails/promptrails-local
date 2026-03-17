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

type CredentialHandler struct {
	store *store.Store
}

func (h *CredentialHandler) List(c echo.Context) error {
	p := getPagination(c)
	creds, total := h.store.ListCredentials(p.Page, p.Limit)
	return listResponse(c, creds, total, p)
}

func (h *CredentialHandler) Create(c echo.Context) error {
	wid := getWorkspaceID()
	var req model.CreateCredentialRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return badRequest(c, "name is required")
	}

	// Mask the value for storage display
	masked := maskValue(req.Value)

	now := time.Now()
	cred := model.Credential{
		ID:            ksuid.New().String(),
		WorkspaceID:   wid,
		Name:          req.Name,
		Type:          req.Type,
		Category:      req.Category,
		Description:   req.Description,
		MaskedContent: masked,
		IsDefault:     req.IsDefault,
		SchemaType:    req.SchemaType,
		IsValid:       true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	h.store.CreateCredential(cred)
	return dataResponse(c, http.StatusCreated, cred)
}

func (h *CredentialHandler) Get(c echo.Context) error {
	cred, ok := h.store.GetCredential(c.Param("credentialId"))
	if !ok {
		return notFound(c, "credential not found")
	}
	return dataResponse(c, http.StatusOK, cred)
}

func (h *CredentialHandler) Update(c echo.Context) error {
	cred, ok := h.store.GetCredential(c.Param("credentialId"))
	if !ok {
		return notFound(c, "credential not found")
	}

	var req model.UpdateCredentialRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if req.Name != nil {
		cred.Name = *req.Name
	}
	if req.Description != nil {
		cred.Description = *req.Description
	}
	if req.Value != nil {
		cred.MaskedContent = maskValue(*req.Value)
	}
	if req.IsDefault != nil {
		cred.IsDefault = *req.IsDefault
	}
	cred.UpdatedAt = time.Now()

	h.store.UpdateCredential(cred)
	return dataResponse(c, http.StatusOK, cred)
}

func (h *CredentialHandler) Delete(c echo.Context) error {
	if !h.store.DeleteCredential(c.Param("credentialId")) {
		return notFound(c, "credential not found")
	}
	return c.NoContent(http.StatusNoContent)
}

// maskValue returns a masked version of the credential value.
func maskValue(val string) string {
	if len(val) <= 8 {
		return strings.Repeat("*", len(val))
	}
	return val[:4] + strings.Repeat("*", len(val)-8) + val[len(val)-4:]
}
