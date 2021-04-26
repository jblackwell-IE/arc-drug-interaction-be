package fdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IdeaEvolver/cutter-pkg/client"
)

type DrugInteractionsResponse struct {
	Response DDIScreenResults
}

type DDIScreenResults struct {
	ScreenResults []ScreenResult
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

func (c *Client) GetDrugInteractions(ctx context.Context, drugIds []string) (*ScreenResult, error) {
	//"SHAREDKEY"+" " + clientid  + ":" + secret
	fmt.Println("Drugs ids", drugIds)
	authString := c.AuthScheme + " " + c.ClientId + ":" + c.Secret
	//TODO write body; nil for now
	req, _ := client.NewRequestWithContext(ctx, "POST", c.FDBUrl, nil)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", authString)

	ret := DrugInteractionsResponse{}
	if err := c.do(ctx, req, &ret); err != nil {
		return nil, err
	}

	return &ret.Response.ScreenResults[0], nil
}
