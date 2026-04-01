package cli

import "github.com/spf13/cobra"

func NewUpdateCommand() *cobra.Command {
	var provider string
	var recordType string
	var value string
	var ttl int

	cmd := &cobra.Command{
		Use:   "update <domain> <name>",
		Short: "Update an existing DNS record through a provider",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVar(&provider, "provider", "", "DNS provider to target")
	cmd.Flags().StringVar(&recordType, "type", "", "Record type to update")
	cmd.Flags().StringVar(&value, "value", "", "New record value")
	cmd.Flags().IntVar(&ttl, "ttl", 600, "New record TTL in seconds")

	return cmd
}
