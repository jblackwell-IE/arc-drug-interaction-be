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
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	FDBEndpoint string `envconfig:"FDB_ENDPOINT" required:"true"`
	AuthScheme  string `envconfig:"AUTH_SCHEME" required:"true"`
	ClientId    string `envconfig:"CLIENT_ID" required:"true"`
	Secret      string `envconfig:"SECRET" required:"true"`
	Port        string `envconfig:"PORT"`
}

func main() {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		clog.Fatalf("config: %s", err)
	}

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
		Client:     client.New(internalClient),
		FDBUrl:     cfg.FDBEndpoint,
		AuthScheme: cfg.AuthScheme,
		ClientId:   cfg.ClientId,
		Secret:     cfg.Secret,
	}

	handler := &server.Handler{
		Interactions: interactionsClient,
	}

	s := server.New(scfg, handler)
	clog.Infof("listening on %s", s.Addr)
	fmt.Println(s.ListenAndServe())
}
