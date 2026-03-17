package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/store"
)

type TraceHandler struct {
	store *store.Store
}

func (h *TraceHandler) List(c echo.Context) error {
	p := getPagination(c)
	traces, total := h.store.ListTraces(p.Page, p.Limit)
	return listResponse(c, traces, total, p)
}

func (h *TraceHandler) Get(c echo.Context) error {
	trace, ok := h.store.GetTrace(c.Param("traceId"))
	if !ok {
		return notFound(c, "trace not found")
	}
	return dataResponse(c, http.StatusOK, trace)
}
