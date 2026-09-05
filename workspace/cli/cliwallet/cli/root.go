package cli

import "github.com/spf13/cobra"

var Cli = &cobra.Command{
	Use:           "cliwallet",
	Short:         "A simple CLI wallet application",
	SilenceUsage:  true,
	SilenceErrors: true,
}
