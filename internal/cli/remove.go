package cli

import "github.com/spf13/cobra"

func NewRemoveCommand() *cobra.Command {
	var provider string
	var recordType string

	cmd := &cobra.Command{
		Use:   "remove <domain> <name>",
		Short: "Remove a DNS record through a provider",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVar(&provider, "provider", "", "DNS provider to target")
	cmd.Flags().StringVar(&recordType, "type", "", "Record type to remove")

	return cmd
}
