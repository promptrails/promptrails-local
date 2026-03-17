package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/store"
)

type LLMModelHandler struct {
	store *store.Store
}

func (h *LLMModelHandler) List(c echo.Context) error {
	models := h.store.ListLLMModels()
	return dataResponse(c, http.StatusOK, models)
}
