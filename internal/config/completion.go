package config

import (
	"strconv"
	"time"

	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
)

func (c *Config) RegisterCompletions(cmd *cobra.Command) {
	must.Must(cmd.RegisterFlagCompletionFunc(FlagBaseURL, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"https://"}, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagYear, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return completeRange(2013, time.Now().Year())
	}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCropLeft, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCropRight, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCropTop, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCropBottom, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagZoom, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return completeRange(2, 6)
	}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagFormat, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"png", "png8", "png24"}, cobra.ShellCompDirectiveNoFileComp
	}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCompression, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return CompressionLevelStrings(), cobra.ShellCompDirectiveNoFileComp
	}))
}

func completeRange(first, last int) ([]string, cobra.ShellCompDirective) {
	s := make([]string, 0, last-first)
	for i := first; i <= last; i++ {
		s = append(s, strconv.Itoa(i))
	}
	return s, cobra.ShellCompDirectiveNoFileComp
}
