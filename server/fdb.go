package server

import (
	"encoding/json"
	"net/http"
)

type DrugInteractionsRequest struct {
	DrugIds []string `json:"drugIds"`
}

func (h *Handler) GetDrugInteractions(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	req := DrugInteractionsRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	return h.FDB.GetDrugInteractions(r.Context(), req.DrugIds)
}
