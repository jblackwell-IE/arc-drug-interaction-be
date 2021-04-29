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

type DrugsRequest struct {
	DrugNames []string `json:"drugNames"`
}

func (h *Handler) CheckDrugInteractions(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	req := DrugInteractionsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return h.checkDrugInteractions(r.Context(), req.DrugIds)
}

func (h *Handler) checkDrugInteractions(ctx context.Context, drugIds []string) (*fdb.DrugInteractionsResponse, error) {
	return h.Interactions.CheckDrugInteractions(ctx, drugIds)
}

func (h *Handler) GetDrugIds(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	req := DrugsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return h.getDrugIds(r.Context(), req.DrugNames)
}

func (h *Handler) getDrugIds(ctx context.Context, drugNames []string) ([]*fdb.DrugIdsResponse, error) {
	fmt.Println("drugs", drugNames)
	return h.Interactions.GetDrugIds(ctx, drugNames)
}
