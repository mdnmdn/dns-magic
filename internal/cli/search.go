package cli

import "github.com/spf13/cobra"

func NewSearchCommand() *cobra.Command {
	var provider string
	var recordType string

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search DNS records across a zone or provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVar(&provider, "provider", "", "DNS provider to search")
	cmd.Flags().StringVar(&recordType, "type", "", "Record type filter")

	return cmd
}
