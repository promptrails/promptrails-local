package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/promptrails/promptrails-local/internal/model"
	"github.com/promptrails/promptrails-local/internal/store"
)

type LLMModelHandler struct {
	store *store.Store
}

func (h *LLMModelHandler) List(c echo.Context) error {
	models := h.store.ListLLMModels()
	return dataResponse(c, http.StatusOK, models)
}

// ListAvailable returns active models grouped by provider, mirroring the real
// available-models endpoint. The emulator does not filter by credentials.
func (h *LLMModelHandler) ListAvailable(c echo.Context) error {
	byProvider := map[string][]model.AvailableModelEntry{}
	order := []string{}
	for _, m := range h.store.ListLLMModels() {
		if !m.IsActive {
			continue
		}
		if _, seen := byProvider[m.Provider]; !seen {
			order = append(order, m.Provider)
		}
		byProvider[m.Provider] = append(byProvider[m.Provider], model.AvailableModelEntry{
			ID:                    m.ID,
			ModelID:               m.ModelID,
			DisplayName:           m.DisplayName,
			MaxTokens:             m.MaxTokens,
			SupportsVision:        m.SupportsVision,
			SupportsTools:         m.SupportsTools,
			SupportsJSON:          m.SupportsJSON,
			SupportsTemperature:   m.SupportsTemperature,
			SupportsTopP:          m.SupportsTopP,
			SupportsTopK:          m.SupportsTopK,
			SupportsReasoning:     m.SupportsReasoning,
			SupportsWebSearch:     m.SupportsWebSearch,
			SupportsPromptCaching: m.SupportsPromptCaching,
			InputPrice:            m.InputPrice,
			OutputPrice:           m.OutputPrice,
			IsDeprecated:          m.IsDeprecated,
		})
	}

	groups := make([]model.AvailableModelGroup, 0, len(order))
	for _, provider := range order {
		groups = append(groups, model.AvailableModelGroup{Provider: provider, Models: byProvider[provider]})
	}
	return dataResponse(c, http.StatusOK, map[string]any{"groups": groups})
}
