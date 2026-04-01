package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	runtime := NewRuntime()

	cmd := &cobra.Command{
		Use:   "dns-magic",
		Short: "Swiss-army knife CLI for DNS inspection and provider operations",
	}

	cmd.PersistentFlags().StringVar(&runtime.ConfigPath, "config", runtime.ConfigPath, "Path to the TOML config file")

	cmd.AddCommand(
		NewCheckCommand(),
		NewListCommand(runtime),
		NewListDomainCommand(runtime),
		NewSearchCommand(),
		NewAddCommand(),
		NewUpdateCommand(),
		NewRemoveCommand(),
	)

	return cmd
}
