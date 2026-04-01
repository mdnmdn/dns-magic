package app

import "github.com/mdnmdn/dns-magic/internal/cli"

func Run() error {
	return cli.NewRootCommand().Execute()
}
