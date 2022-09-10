package cmd

import (
	"github.com/urfave/cli/v2"
)

func init() {
	registerCommand(serverCmd)
}

var serverCmd = &cli.Command{
	Name:    "server",
	Aliases: []string{"srv"},
	Usage:   "start a server that receives files and listens on a specified port.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "network",
			Aliases: []string{"nw"},
			Usage:   "choose a network protocol(tcp|udp)",
			Value:   "tcp",
		},
		&cli.StringFlag{
			Name:  "addr",
			Usage: "specify a listening port",
			Value: "9988",
		},
	},
	Action: func(ctx *cli.Context) error {

		return nil
	},
}
