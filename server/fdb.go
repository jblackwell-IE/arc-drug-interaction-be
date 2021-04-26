package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jblackwell-IE/arc-drug-interaction-be/fdb"
)

type DrugInteractionsRequest struct {
	DrugIds []string `json:"drugIds"`
}

func (h *Handler) GetDrugInteractions(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	req := DrugInteractionsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	fmt.Println("req", req)

	return h.getDrugInteractions(r.Context(), req.DrugIds)
}

func (h *Handler) getDrugInteractions(ctx context.Context, drugIds []string) (*fdb.ScreenResult, error) {
	return h.Interactions.GetDrugInteractions(ctx, drugIds)
}
