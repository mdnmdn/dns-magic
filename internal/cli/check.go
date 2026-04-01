package cli

import "github.com/spf13/cobra"

func NewCheckCommand() *cobra.Command {
	var dnsIP string

	cmd := &cobra.Command{
		Use:   "check <name>",
		Short: "Check DNS records for a name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVar(&dnsIP, "dns", "", "Optional DNS server IP to query")

	return cmd
}
