package cli

import "github.com/spf13/cobra"

func NewAddCommand() *cobra.Command {
	var provider string
	var recordType string
	var value string
	var ttl int

	cmd := &cobra.Command{
		Use:   "add <domain> <name>",
		Short: "Add a DNS record through a provider",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVar(&provider, "provider", "", "DNS provider to target")
	cmd.Flags().StringVar(&recordType, "type", "", "Record type to create")
	cmd.Flags().StringVar(&value, "value", "", "Record value")
	cmd.Flags().IntVar(&ttl, "ttl", 600, "Record TTL in seconds")

	return cmd
}
