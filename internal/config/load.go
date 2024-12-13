package config

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"gabe565.com/utils/cobrax"
	"gabe565.com/utils/httpx"
	"github.com/spf13/cobra"
)

var ErrMissingConfig = errors.New("command missing config")

func Load(ctx context.Context, cmd *cobra.Command) (*Config, error) {
	conf, ok := FromContext(cmd.Context())
	if !ok {
		return conf, ErrMissingConfig
	}

	if err := conf.DetermineOffsetsByYear(); err != nil {
		return conf, err
	}

	if conf.Completion == "" {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: conf.Insecure} //nolint:gosec
		conf.Client.Transport = httpx.NewUserAgentTransport(transport, cobrax.BuildUserAgent(cmd))

		if err := conf.CheckYear(ctx); err != nil {
			return conf, err
		}

		if err := conf.FindFormat(ctx); err != nil {
			return conf, err
		}
	}

	return conf, nil
}
