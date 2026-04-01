package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/mdnmdn/dns-magic/internal/output"
	"github.com/mdnmdn/dns-magic/internal/providers"
)

func NewListCommand(runtime *Runtime) *cobra.Command {
	var recordType string
	var provider string
	var shopperID string
	var outputFormat string
	var name string
	var offset int
	var limit int

	cmd := &cobra.Command{
		Use:   "list <domain>",
		Short: "List DNS records for a domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOutput(outputFormat); err != nil {
				return err
			}

			client, _, err := runtime.MustProvider(context.Background(), provider)
			if err != nil {
				return err
			}

			records, err := client.ListRecords(
				cmd.Context(),
				providers.RecordListOptions{
					Domain:    args[0],
					Type:      recordType,
					Name:      name,
					ShopperID: shopperID,
					Offset:    offset,
					Limit:     limit,
				},
			)
			if err != nil {
				return err
			}

			return output.Write(os.Stdout, outputFormat, records)
		},
	}

	cmd.Flags().StringVar(&recordType, "type", "", "Record type filter, for example A or CNAME")
	cmd.Flags().StringVar(&provider, "provider", "", "Provider alias from config, for example godaddy:customer1")
	cmd.Flags().StringVar(&shopperID, "shopper-id", "", "GoDaddy delegated shopper override")
	cmd.Flags().StringVar(&outputFormat, "output", output.FormatTable, "Output format: table, json, yaml, markdown, toon")
	cmd.Flags().StringVar(&name, "name", "", "Record name filter; requires --type")
	cmd.Flags().IntVar(&offset, "offset", 0, "Pagination offset for provider APIs that support it")
	cmd.Flags().IntVar(&limit, "limit", 0, "Maximum number of records to return")

	return cmd
}
