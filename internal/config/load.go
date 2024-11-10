package config

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
)

var ErrMissingConfig = errors.New("command missing config")

func Load(ctx context.Context, cmd *cobra.Command) (*Config, error) {
	InitLog()

	conf, ok := FromContext(cmd.Context())
	if !ok {
		return conf, ErrMissingConfig
	}

	if err := conf.DetermineOffsetsByYear(); err != nil {
		return conf, err
	}

	if conf.Completion == "" {
		if err := conf.CheckYear(ctx); err != nil {
			return conf, err
		}

		if err := conf.FindFormat(ctx); err != nil {
			return conf, err
		}
	}

	return conf, nil
}
