package cli

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mdnmdn/dns-magic/internal/output"
	"github.com/mdnmdn/dns-magic/internal/providers"
)

func NewListDomainCommand(runtime *Runtime) *cobra.Command {
	var provider string
	var shopperID string
	var outputFormat string
	var marker string
	var modifiedSince string
	var limit int
	var statuses []string
	var statusGroups []string
	var includes []string

	cmd := &cobra.Command{
		Use:   "list-domain",
		Short: "List domains available through a provider account",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOutput(outputFormat); err != nil {
				return err
			}

			client, _, err := runtime.MustProvider(context.Background(), provider)
			if err != nil {
				return err
			}

			domains, err := client.ListDomains(
				cmd.Context(),
				providers.DomainListOptions{
					ShopperID:     shopperID,
					Statuses:      splitAndTrim(statuses),
					StatusGroups:  splitAndTrim(statusGroups),
					Includes:      splitAndTrim(includes),
					Limit:         limit,
					Marker:        marker,
					ModifiedSince: modifiedSince,
				},
			)
			if err != nil {
				return err
			}

			return output.Write(os.Stdout, outputFormat, domains)
		},
	}

	cmd.Flags().StringVar(&provider, "provider", "", "Provider alias from config, for example godaddy:customer1")
	cmd.Flags().StringVar(&shopperID, "shopper-id", "", "GoDaddy delegated shopper override")
	cmd.Flags().StringVar(&outputFormat, "output", output.FormatTable, "Output format: table, json, yaml, markdown, toon")
	cmd.Flags().StringVar(&marker, "marker", "", "Marker domain to continue listing from")
	cmd.Flags().StringVar(&modifiedSince, "modified-since", "", "Only include domains modified since the given RFC3339 timestamp")
	cmd.Flags().IntVar(&limit, "limit", 0, "Maximum number of domains to return")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "Domain status filter; can be provided multiple times or as a comma-separated list")
	cmd.Flags().StringSliceVar(&statusGroups, "status-group", nil, "Domain status group filter; can be provided multiple times or as a comma-separated list")
	cmd.Flags().StringSliceVar(&includes, "include", nil, "Optional response detail, for example nameServers; can be provided multiple times or as a comma-separated list")

	return cmd
}

func splitAndTrim(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			item = strings.TrimSpace(item)
			if item != "" {
				out = append(out, item)
			}
		}
	}

	return out
}
