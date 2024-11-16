package config

import (
	"errors"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func (c *Config) RegisterCompletions(cmd *cobra.Command) {
	if err := errors.Join(
		cmd.RegisterFlagCompletionFunc(FlagYear, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return completeRange(2013, time.Now().Year())
		}),
		cmd.RegisterFlagCompletionFunc(FlagTileMinX, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return completeRange(0, 63)
		}),
		cmd.RegisterFlagCompletionFunc(FlagTileMaxX, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return completeRange(0, 63)
		}),
		cmd.RegisterFlagCompletionFunc(FlagTileMinY, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return completeRange(0, 63)
		}),
		cmd.RegisterFlagCompletionFunc(FlagTileMaxY, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return completeRange(0, 63)
		}),
		cmd.RegisterFlagCompletionFunc(FlagZoom, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return completeRange(2, 6)
		}),
		cmd.RegisterFlagCompletionFunc(FlagFormat, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"png", "png8", "png24"}, cobra.ShellCompDirectiveNoFileComp
		}),
		cmd.RegisterFlagCompletionFunc(FlagCompression, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return CompressionLevelStrings(), cobra.ShellCompDirectiveNoFileComp
		}),
	); err != nil {
		panic(err)
	}
}

func completeRange(first, last int) ([]string, cobra.ShellCompDirective) {
	s := make([]string, 0, last-first)
	for i := first; i <= last; i++ {
		s = append(s, strconv.Itoa(i))
	}
	return s, cobra.ShellCompDirectiveNoFileComp
}
