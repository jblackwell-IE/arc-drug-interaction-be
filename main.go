package main

import "github.com/kelseyhightower/envconfig"

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
}
