package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/IdeaEvolver/cutter-pkg/client"
	"github.com/IdeaEvolver/cutter-pkg/service"
	"go.opencensus.io/plugin/ochttp"

	"github.com/IdeaEvolver/cutter-pkg/clog"
	"github.com/jblackwell-IE/arc-drug-interaction-be/fdb"
	"github.com/jblackwell-IE/arc-drug-interaction-be/server"
)

type Config struct {
	InteractionsEndpoint string `envconfig:"FDB_INTERACTIONS" required:"true"`
	DrugIdsEndpoint      string `envconfig:"FDB_DRUG_IDS" required:"true"`
	AuthScheme           string `envconfig:"AUTH_SCHEME" required:"true"`
	ClientId             string `envconfig:"CLIENT_ID" required:"true"`
	Secret               string `envconfig:"SECRET" required:"true"`
	Port                 string `envconfig:"PORT"`
}

func main() {
	cfg := &Config{}
	// if err := envconfig.Process("", cfg); err != nil {
	// 	clog.Fatalf("config: %s", err)
	// }

	cfg.InteractionsEndpoint = "https://api.fdbcloudconnector.com/CC/api/v1_4/Screen"
	cfg.DrugIdsEndpoint = "https://api.fdbcloudconnector.com/CC/api/v1_4/PrescribableDrugs"
	cfg.AuthScheme = "SHAREDKEY"
	cfg.ClientId = "1777"
	cfg.Secret = "x/RMaGKqBE8KUX8o4qM/V3ZsenNfE6S0ZSQBrV74PM4="
	cfg.Port = "8080"

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	internalClient := &http.Client{
		Transport: &ochttp.Transport{
			// Use Google Cloud propagation format.
			Propagation: &propagation.HTTPFormat{},
			Base:        customTransport,
		},
	}

	scfg := &service.Config{
		Addr:                fmt.Sprintf(":%s", cfg.Port),
		ShutdownGracePeriod: time.Second * 10,
		MaxShutdownTime:     time.Second * 30,
	}

	interactionsClient := &fdb.Client{
		Client:               client.New(internalClient),
		InteractionsEndpoint: cfg.InteractionsEndpoint,
		DrugIdsEndpoint:      cfg.DrugIdsEndpoint,
		Auth:                 cfg.AuthScheme + " " + cfg.ClientId + ":" + cfg.Secret,
	}

	handler := &server.Handler{
		Interactions: interactionsClient,
	}

	s := server.New(scfg, handler)
	clog.Infof("listening on %s", s.Addr)
	fmt.Println(s.ListenAndServe())
}
