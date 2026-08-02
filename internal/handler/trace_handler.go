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

// Summary returns aggregate metering statistics over all traces — the v2
// replacement for the removed cost summary. Query filters are accepted for
// parity with the real API but the single-namespace emulator aggregates all.
func (h *TraceHandler) Summary(c echo.Context) error {
	return dataResponse(c, http.StatusOK, h.store.TraceSummary())
}
