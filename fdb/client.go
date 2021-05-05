package fdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IdeaEvolver/cutter-pkg/client"
)

type DrugInteractionsResponse struct {
	DrugId   string `json:"drugId"`
	Severity string `json:"severity"`
}

type Interactions struct {
	DDIScreenResponse DDIScreenResponse `json:"DDIScreenResponse"`
}

type DDIScreenResponse struct {
	DDIScreenResults []ScreenResult `json:"DDIScreenResults"`
}

type ScreenResult struct {
	Severity string `json:"Severity"`
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

type DrugIdsResponse struct {
	DrugName string `json:"drugName"`
	DrugId   string `json:"drugId"`
}

type ItemResults struct {
	Items []DrugResult `json:"Items"`
}

type DrugResult struct {
	PrescribableDrugID string `json:"PrescribableDrugID"`
}

type Client struct {
	Client               *client.Client
	InteractionsEndpoint string
	DrugIdsEndpoint      string
	Auth                 string
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

var amlodipineId = "151400"

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

func (c *Client) CheckDrugInteractions(ctx context.Context, drugIds []string) ([]*DrugInteractionsResponse, error) {
	ret := []*DrugInteractionsResponse{}

	for _, id := range drugIds {
		drugsReq := &DrugInteractionsRequest{
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
						DrugId:          amlodipineId,
						DrugDesc:        nil,
						DrugConceptType: "3", // field is required TODO how is this number calculated? Hardcoded for now.
					},
					ScreenDrug{
						Prospective:     false,
						DrugId:          id,
						DrugDesc:        nil,
						DrugConceptType: "3", // same TODO
					},
				},
			},
		}

		b, _ := json.Marshal(drugsReq)

		req, _ := client.NewRequestWithContext(ctx, "POST", c.InteractionsEndpoint, bytes.NewReader(b))
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", c.Auth)

		interactions := Interactions{}
		if err := c.do(ctx, req, &interactions); err != nil {
			return nil, err
		}

		resp := DrugInteractionsResponse{}
		resp.DrugId = id

		if len(interactions.DDIScreenResponse.DDIScreenResults) == 0 {
			resp.Severity = ""
		} else {
			resp.Severity = interactions.DDIScreenResponse.DDIScreenResults[0].Severity
		}

		ret = append(ret, &resp)

	}

	return ret, nil
}

func (c *Client) GetDrugIds(ctx context.Context, drugNames []string) ([]*DrugIdsResponse, error) {
	drugIds := []*DrugIdsResponse{}
	for _, name := range drugNames {
		url := fmt.Sprintf(c.DrugIdsEndpoint+"?callSystemName=test&callid=123&searchtext=%s&searchtype=startswith", name)
		req, _ := client.NewRequestWithContext(ctx, "GET", url, nil)
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", c.Auth)

		items := ItemResults{}
		if err := c.do(ctx, req, &items); err != nil {
			return nil, err
		}
		drugId := DrugIdsResponse{}
		drugId.DrugName = name
		drugId.DrugId = items.Items[0].PrescribableDrugID

		drugIds = append(drugIds, &drugId)
	}

	return drugIds, nil
}

//https://api.fdbcloudconnector.com/CC/api/v1_4/PrescribableDrugs?callSystemName=test&callid=123&searchtext=crestor&searchtype=startswith
