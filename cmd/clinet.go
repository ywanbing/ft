package cmd

import (
	"github.com/urfave/cli/v2"
)

var clientCmd = &cli.Command{
	Name:    "client",
	Aliases: []string{"cli"},
	Usage:   "start an upload client.",
	Flags:   []cli.Flag{},
	Action: func(ctx *cli.Context) error {

		return nil
	},
}

func init() {
	registerCommand(clientCmd)
}
