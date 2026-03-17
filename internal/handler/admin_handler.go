package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/seed"
	"github.com/promptrails/promptrails-local/internal/store"
	"go.uber.org/zap"
)

type AdminHandler struct {
	store  *store.Store
	logger *zap.Logger
}

// Reset clears all data and reloads seed data.
func (h *AdminHandler) Reset(c echo.Context) error {
	h.store.Reset()

	if err := seed.Load(h.store, h.logger); err != nil {
		h.logger.Error("failed to reload seed data after reset", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "reset succeeded but seed reload failed: " + err.Error(),
		})
	}

	h.logger.Info("store reset and seed data reloaded")
	return dataResponse(c, http.StatusOK, map[string]string{"status": "reset complete"})
}

// Seed reloads seed data without clearing existing data.
func (h *AdminHandler) Seed(c echo.Context) error {
	if err := seed.Load(h.store, h.logger); err != nil {
		h.logger.Error("failed to reload seed data", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "seed reload failed: " + err.Error(),
		})
	}

	h.logger.Info("seed data reloaded")
	return dataResponse(c, http.StatusOK, map[string]string{"status": "seed complete"})
}

// Stats returns counts for every entity type in the store.
func (h *AdminHandler) Stats(c echo.Context) error {
	stats := h.store.Stats()
	return dataResponse(c, http.StatusOK, stats)
}
