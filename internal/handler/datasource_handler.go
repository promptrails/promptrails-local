package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
	"github.com/segmentio/ksuid"
)

type DataSourceHandler struct {
	store *store.Store
}

func (h *DataSourceHandler) List(c echo.Context) error {
	p := getPagination(c)
	sources, total := h.store.ListDataSources(p.Page, p.Limit)
	return listResponse(c, sources, total, p)
}

func (h *DataSourceHandler) Create(c echo.Context) error {
	wid := getWorkspaceID()
	var req model.CreateDataSourceRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return badRequest(c, "name is required")
	}

	now := time.Now()
	ds := model.DataSource{
		ID:          ksuid.New().String(),
		WorkspaceID: wid,
		Name:        req.Name,
		Type:        req.Type,
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.store.CreateDataSource(ds)
	return dataResponse(c, http.StatusCreated, ds)
}

func (h *DataSourceHandler) Get(c echo.Context) error {
	ds, ok := h.store.GetDataSource(c.Param("dataSourceId"))
	if !ok {
		return notFound(c, "data source not found")
	}
	return dataResponse(c, http.StatusOK, ds)
}

func (h *DataSourceHandler) Update(c echo.Context) error {
	ds, ok := h.store.GetDataSource(c.Param("dataSourceId"))
	if !ok {
		return notFound(c, "data source not found")
	}

	var req model.UpdateDataSourceRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if req.Name != nil {
		ds.Name = *req.Name
	}
	if req.Description != nil {
		ds.Description = *req.Description
	}
	if req.Status != nil {
		ds.Status = *req.Status
	}
	ds.UpdatedAt = time.Now()

	h.store.UpdateDataSource(ds)
	return dataResponse(c, http.StatusOK, ds)
}

func (h *DataSourceHandler) Delete(c echo.Context) error {
	if !h.store.DeleteDataSource(c.Param("dataSourceId")) {
		return notFound(c, "data source not found")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *DataSourceHandler) ListVersions(c echo.Context) error {
	dsID := c.Param("dataSourceId")
	if _, ok := h.store.GetDataSource(dsID); !ok {
		return notFound(c, "data source not found")
	}
	versions := h.store.ListDataSourceVersions(dsID)
	return dataResponse(c, http.StatusOK, versions)
}

func (h *DataSourceHandler) CreateVersion(c echo.Context) error {
	dsID := c.Param("dataSourceId")
	if _, ok := h.store.GetDataSource(dsID); !ok {
		return notFound(c, "data source not found")
	}

	var req model.CreateDataSourceVersionRequest
	if err := c.Bind(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	version := model.DataSourceVersion{
		ID:               ksuid.New().String(),
		DataSourceID:     dsID,
		Version:          req.Version,
		CredentialID:     req.CredentialID,
		ConnectionConfig: req.ConnectionConfig,
		QueryTemplate:    req.QueryTemplate,
		Parameters:       req.Parameters,
		IsCurrent:        req.SetCurrent,
		Message:          req.Message,
		CacheTimeout:     req.CacheTimeout,
		OutputFormat:     req.OutputFormat,
		CreatedAt:        time.Now(),
	}

	if req.SetCurrent {
		h.store.DemoteDataSourceVersions(dsID)
	}

	h.store.CreateDataSourceVersion(version)
	return dataResponse(c, http.StatusCreated, version)
}

func (h *DataSourceHandler) PromoteVersion(c echo.Context) error {
	dsID := c.Param("dataSourceId")
	versionID := c.Param("versionId")

	version, ok := h.store.GetDataSourceVersion(versionID)
	if !ok || version.DataSourceID != dsID {
		return notFound(c, "version not found")
	}

	h.store.DemoteDataSourceVersions(dsID)
	version.IsCurrent = true
	h.store.UpdateDataSourceVersion(version)
	return dataResponse(c, http.StatusOK, version)
}

func (h *DataSourceHandler) Query(c echo.Context) error {
	dsID := c.Param("dataSourceId")
	if _, ok := h.store.GetDataSource(dsID); !ok {
		return notFound(c, "data source not found")
	}

	// Return mock query results
	mockResults := map[string]any{
		"columns": []string{"id", "name", "value"},
		"rows": [][]any{
			{1, "sample_row_1", 42},
			{2, "sample_row_2", 99},
		},
		"row_count":   2,
		"duration_ms": 12,
	}
	return dataResponse(c, http.StatusOK, mockResults)
}
