package handler

import (
	"errors"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"github.com/segmentio/ksuid"
)

// AgentVFSHandler serves the per-agent Virtual Filesystem in the emulator.
//
// The emulator's VFS is an in-memory mock with the same shape as the platform
// but no quota enforcement and no MIME sniffing — enough for SDK tests and
// local smoke runs.
type AgentVFSHandler struct {
	store *store.Store
}

const (
	vfsMaxNameLen = 255
	vfsMaxDepth   = 16
)

var errVFSInvalidPath = errors.New("invalid vfs path")

// normalizeVFSPath validates and canonicalizes a path. Same rules as the
// platform: absolute, no traversal, no NUL, depth ≤ 16, names ≤ 255.
func normalizeVFSPath(p string) (string, error) {
	if p == "" {
		return "", errVFSInvalidPath
	}
	if strings.ContainsRune(p, 0x00) {
		return "", errVFSInvalidPath
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	for _, seg := range strings.Split(p, "/") {
		if seg == ".." {
			return "", errVFSInvalidPath
		}
	}
	cleaned := path.Clean(p)
	if cleaned == "/" {
		return "/", nil
	}
	segments := strings.Split(strings.TrimPrefix(cleaned, "/"), "/")
	if len(segments) > vfsMaxDepth {
		return "", errVFSInvalidPath
	}
	for _, seg := range segments {
		if seg == "" || seg == "." || seg == ".." {
			return "", errVFSInvalidPath
		}
		if len(seg) > vfsMaxNameLen {
			return "", errVFSInvalidPath
		}
	}
	return cleaned, nil
}

func splitVFSPath(p string) (parent, name string) {
	if p == "/" || p == "" {
		return "", ""
	}
	idx := strings.LastIndex(p, "/")
	if idx <= 0 {
		return "/", strings.TrimPrefix(p, "/")
	}
	return p[:idx], p[idx+1:]
}

func (h *AgentVFSHandler) ensureDirs(agentID, target string) {
	if target == "/" || target == "" {
		return
	}
	segments := strings.Split(strings.TrimPrefix(target, "/"), "/")
	current := ""
	for _, seg := range segments {
		current = current + "/" + seg
		if _, ok := h.store.VFSGet(agentID, current); ok {
			continue
		}
		parent, name := splitVFSPath(current)
		now := time.Now()
		h.store.VFSUpsert(model.AgentVFSFile{
			ID:             ksuid.New().String(),
			AgentID:        agentID,
			Path:           current,
			ParentPath:     parent,
			Name:           name,
			IsDir:          true,
			LastWriterKind: "user",
			CreatedAt:      now,
			UpdatedAt:      now,
		})
	}
}

// List returns directory entries.
func (h *AgentVFSHandler) List(c echo.Context) error {
	agentID := c.Param("agentId")
	p := c.QueryParam("path")
	if p == "" {
		p = "/"
	}
	normalized, err := normalizeVFSPath(p)
	if err != nil {
		return badRequest(c, err.Error())
	}
	recursive := c.QueryParam("recursive") == "true"
	items := h.store.VFSList(agentID, normalized, recursive)
	return dataResponse(c, http.StatusOK, map[string]any{
		"path":  normalized,
		"items": items,
		"total": len(items),
	})
}

// Read returns the file content (no line-range support in the emulator;
// the response shape still matches the platform so SDKs work unchanged).
func (h *AgentVFSHandler) Read(c echo.Context) error {
	agentID := c.Param("agentId")
	normalized, err := normalizeVFSPath(c.QueryParam("path"))
	if err != nil {
		return badRequest(c, err.Error())
	}
	file, ok := h.store.VFSGet(agentID, normalized)
	if !ok {
		return notFound(c, "vfs file not found")
	}
	if file.IsDir {
		return badRequest(c, "path is a directory")
	}
	lines := strings.Split(file.Content, "\n")
	return dataResponse(c, http.StatusOK, map[string]any{
		"file":        file,
		"content":     file.Content,
		"total_lines": len(lines),
		"truncated":   false,
	})
}

// Stat returns metadata for a single path.
func (h *AgentVFSHandler) Stat(c echo.Context) error {
	agentID := c.Param("agentId")
	normalized, err := normalizeVFSPath(c.QueryParam("path"))
	if err != nil {
		return badRequest(c, err.Error())
	}
	file, ok := h.store.VFSGet(agentID, normalized)
	if !ok {
		return notFound(c, "vfs file not found")
	}
	return dataResponse(c, http.StatusOK, file)
}

// Write creates or updates a file (overwrite or append).
func (h *AgentVFSHandler) Write(c echo.Context) error {
	agentID := c.Param("agentId")
	var req struct {
		Path     string `json:"path"`
		Content  string `json:"content"`
		Mode     string `json:"mode"`
		MimeType string `json:"mime_type"`
	}
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	normalized, err := normalizeVFSPath(req.Path)
	if err != nil {
		return badRequest(c, err.Error())
	}
	if normalized == "/" {
		return badRequest(c, "cannot write to root")
	}
	if existing, ok := h.store.VFSGet(agentID, normalized); ok && existing.IsDir {
		return badRequest(c, "path is a directory")
	}
	parent, name := splitVFSPath(normalized)
	h.ensureDirs(agentID, parent)
	content := req.Content
	if req.Mode == "append" {
		if existing, ok := h.store.VFSGet(agentID, normalized); ok {
			content = existing.Content + content
		}
	}
	now := time.Now()
	file := model.AgentVFSFile{
		ID:             ksuid.New().String(),
		AgentID:        agentID,
		Path:           normalized,
		ParentPath:     parent,
		Name:           name,
		IsDir:          false,
		Content:        content,
		SizeBytes:      int64(len(content)),
		MimeType:       req.MimeType,
		LastWriterKind: "user",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	h.store.VFSUpsert(file)
	return dataResponse(c, http.StatusOK, file)
}

// Mkdir creates a directory (parents auto-created).
func (h *AgentVFSHandler) Mkdir(c echo.Context) error {
	agentID := c.Param("agentId")
	var req struct {
		Path string `json:"path"`
	}
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	normalized, err := normalizeVFSPath(req.Path)
	if err != nil {
		return badRequest(c, err.Error())
	}
	if normalized == "/" {
		return badRequest(c, "cannot mkdir root")
	}
	h.ensureDirs(agentID, normalized)
	file, ok := h.store.VFSGet(agentID, normalized)
	if !ok {
		return notFound(c, "vfs file not found after mkdir")
	}
	return dataResponse(c, http.StatusCreated, file)
}

// Delete removes a file. recursive=true removes a directory subtree.
func (h *AgentVFSHandler) Delete(c echo.Context) error {
	agentID := c.Param("agentId")
	normalized, err := normalizeVFSPath(c.QueryParam("path"))
	if err != nil {
		return badRequest(c, err.Error())
	}
	if normalized == "/" {
		return badRequest(c, "cannot delete root")
	}
	recursive := c.QueryParam("recursive") == "true"
	deleted := h.store.VFSDelete(agentID, normalized, recursive)
	return dataResponse(c, http.StatusOK, map[string]int{"deleted": deleted})
}

// Move renames or moves a file or directory.
func (h *AgentVFSHandler) Move(c echo.Context) error {
	agentID := c.Param("agentId")
	var req struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	from, err := normalizeVFSPath(req.From)
	if err != nil {
		return badRequest(c, "invalid from: "+err.Error())
	}
	to, err := normalizeVFSPath(req.To)
	if err != nil {
		return badRequest(c, "invalid to: "+err.Error())
	}
	if from == "/" || to == "/" {
		return badRequest(c, "cannot move root")
	}
	if strings.HasPrefix(to, from+"/") {
		return badRequest(c, "cannot move into own descendant")
	}
	src, ok := h.store.VFSGet(agentID, from)
	if !ok {
		return notFound(c, "vfs source not found")
	}
	h.store.VFSDelete(agentID, from, false)
	parent, name := splitVFSPath(to)
	src.Path = to
	src.ParentPath = parent
	src.Name = name
	src.UpdatedAt = time.Now()
	h.store.VFSUpsert(src)
	if src.IsDir {
		for _, child := range h.store.VFSList(agentID, from, true) {
			if child.Path == from {
				continue
			}
			suffix := strings.TrimPrefix(child.Path, from)
			h.store.VFSDelete(agentID, child.Path, false)
			newPath := to + suffix
			p, n := splitVFSPath(newPath)
			child.Path = newPath
			child.ParentPath = p
			child.Name = n
			child.UpdatedAt = time.Now()
			h.store.VFSUpsert(child)
		}
	}
	return c.NoContent(http.StatusNoContent)
}

// Copy duplicates a file or subtree.
func (h *AgentVFSHandler) Copy(c echo.Context) error {
	agentID := c.Param("agentId")
	var req struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	from, err := normalizeVFSPath(req.From)
	if err != nil {
		return badRequest(c, "invalid from: "+err.Error())
	}
	to, err := normalizeVFSPath(req.To)
	if err != nil {
		return badRequest(c, "invalid to: "+err.Error())
	}
	if from == "/" || to == "/" {
		return badRequest(c, "cannot copy root")
	}
	if strings.HasPrefix(to, from+"/") {
		return badRequest(c, "cannot copy into own descendant")
	}
	src, ok := h.store.VFSGet(agentID, from)
	if !ok {
		return notFound(c, "vfs source not found")
	}
	now := time.Now()
	parent, name := splitVFSPath(to)
	dst := src
	dst.ID = ksuid.New().String()
	dst.Path = to
	dst.ParentPath = parent
	dst.Name = name
	dst.CreatedAt = now
	dst.UpdatedAt = now
	h.store.VFSUpsert(dst)
	if src.IsDir {
		for _, child := range h.store.VFSList(agentID, from, true) {
			if child.Path == from {
				continue
			}
			suffix := strings.TrimPrefix(child.Path, from)
			newPath := to + suffix
			p, n := splitVFSPath(newPath)
			copy := child
			copy.ID = ksuid.New().String()
			copy.Path = newPath
			copy.ParentPath = p
			copy.Name = n
			copy.CreatedAt = now
			copy.UpdatedAt = now
			h.store.VFSUpsert(copy)
		}
	}
	return c.NoContent(http.StatusNoContent)
}

// Grep is a case-insensitive content search across files.
func (h *AgentVFSHandler) Grep(c echo.Context) error {
	agentID := c.Param("agentId")
	query := c.QueryParam("q")
	if query == "" {
		return badRequest(c, "q is required")
	}
	within := c.QueryParam("path")
	limit := 100
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	lowered := strings.ToLower(query)
	matches := []map[string]any{}
	for _, f := range h.store.VFSList(agentID, "/", true) {
		if f.IsDir {
			continue
		}
		if within != "" && within != "/" && (f.Path != within && !strings.HasPrefix(f.Path, within+"/")) {
			continue
		}
		for i, line := range strings.Split(f.Content, "\n") {
			if strings.Contains(strings.ToLower(line), lowered) {
				matches = append(matches, map[string]any{
					"path":        f.Path,
					"line_number": i + 1,
					"line":        line,
				})
				if len(matches) >= limit {
					return dataResponse(c, http.StatusOK, map[string]any{"query": query, "matches": matches})
				}
			}
		}
	}
	return dataResponse(c, http.StatusOK, map[string]any{"query": query, "matches": matches})
}

// Glob finds files by name or path pattern. * matches anything, ? matches one char.
func (h *AgentVFSHandler) Glob(c echo.Context) error {
	agentID := c.Param("agentId")
	pattern := c.QueryParam("pattern")
	if pattern == "" {
		return badRequest(c, "pattern is required")
	}
	within := c.QueryParam("path")
	hasSlash := strings.Contains(pattern, "/")
	items := []model.AgentVFSFile{}
	for _, f := range h.store.VFSList(agentID, "/", true) {
		if within != "" && within != "/" && (f.Path != within && !strings.HasPrefix(f.Path, within+"/")) {
			continue
		}
		target := f.Name
		if hasSlash {
			target = f.Path
		}
		if globMatch(pattern, target) {
			items = append(items, f)
		}
	}
	return dataResponse(c, http.StatusOK, map[string]any{"pattern": pattern, "items": items})
}

// Usage returns total bytes used by the agent's VFS.
func (h *AgentVFSHandler) Usage(c echo.Context) error {
	agentID := c.Param("agentId")
	return dataResponse(c, http.StatusOK, map[string]int64{"bytes_used": h.store.VFSUsageBytes(agentID)})
}

// globMatch is a minimal * / ? matcher — sufficient for emulator usage.
func globMatch(pattern, s string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
		return strings.Contains(s, strings.Trim(pattern, "*"))
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(s, strings.TrimPrefix(pattern, "*"))
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(s, strings.TrimSuffix(pattern, "*"))
	}
	return s == pattern
}
