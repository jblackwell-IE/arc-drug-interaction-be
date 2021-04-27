package fdb

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/IdeaEvolver/cutter-pkg/client"
)

type DrugInteractionsResponse struct {
	DDIScreenResponse DDIScreenResponse `json:"DDIScreenResponse"`
}

type DDIScreenResponse struct {
	DDIScreenResults []ScreenResult `json:"DDIScreenResults"`
}

type ScreenResult struct {
	Severity string `json:"Severity"`
}

type Client struct {
	Client     *client.Client
	FDBUrl     string
	AuthScheme string
	ClientId   string
	Secret     string
}

type DrugInteractionsRequest struct {
	DDiscreenRequest struct {
		SeverityFilter int `json:"severityFilter"`
	} `json:"ddiscreenRequest"`
	CallContext struct {
		CallSystemName string `json:"callSystemName"`
	} `json:"callContext"`
	ScreenProfile ScreenProfile `json:"screenProfile"`
}

type ScreenProfile struct {
	ScreenDrugs []ScreenDrug `json:"screenDrugs"`
}

type ScreenDrug struct {
	Prospective     bool    `json:"prospective"`
	DrugId          string  `json:"drugID"`
	DrugDesc        *string `json:"drugDesc"`
	DrugConceptType string  `json:"DrugConceptType"`
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func (c *Client) do(ctx context.Context, req *client.Request, ret interface{}) error {
	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if ret != nil {
		return json.NewDecoder(res.Body).Decode(&ret)
	}

	return nil
}

func (c *Client) GetDrugInteractions(ctx context.Context, drugIds []string) (*DrugInteractionsResponse, error) {
	authString := c.AuthScheme + " " + c.ClientId + ":" + c.Secret

	interactions := &DrugInteractionsRequest{
		DDiscreenRequest: struct {
			SeverityFilter int `json:"severityFilter"`
		}{
			SeverityFilter: 9,
		},
		CallContext: struct {
			CallSystemName string `json:"callSystemName"`
		}{
			CallSystemName: "Test",
		},
		ScreenProfile: ScreenProfile{
			ScreenDrugs: []ScreenDrug{
				ScreenDrug{
					Prospective:     false, // TODO figure out this field, left as false for now based of FDB docs
					DrugId:          drugIds[0],
					DrugDesc:        nil,
					DrugConceptType: "2", // field is required TODO how is this number calculated? Hardcoded for now.
				},
				ScreenDrug{
					Prospective:     false,
					DrugId:          drugIds[1],
					DrugDesc:        nil,
					DrugConceptType: "3", // same TODO
				},
			},
		},
	}

	b, _ := json.Marshal(interactions)

	req, _ := client.NewRequestWithContext(ctx, "POST", c.FDBUrl, bytes.NewReader(b))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", authString)

	ret := DrugInteractionsResponse{}
	if err := c.do(ctx, req, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}
