package cmd

import (
	"fmt"
	"net"

	"github.com/urfave/cli/v2"
	"github.com/ywanbing/ft/pkg/file"
	"github.com/ywanbing/ft/pkg/server"
)

var clientCmd = &cli.Command{
	Name:    "client",
	Aliases: []string{"cli"},
	Usage:   "start an upload client.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "network",
			Aliases: []string{"nw"},
			Usage:   "choose a network protocol(tcp|udp)",
			Value:   "tcp",
		},
		&cli.StringFlag{
			Name:  "addr",
			Usage: "specify a server address",
			Value: "127.0.0.1:9988",
		},
		&cli.StringFlag{
			Name:    "dir",
			Value:   "./",
			Aliases: []string{"d"},
			Usage:   "upload dir or save dir",
		},
	},
	Action: func(ctx *cli.Context) error {
		network := ctx.String("network")
		addr := ctx.String("addr")
		dir := ctx.String("dir")

		if !file.PathExists(dir) {
			return fmt.Errorf("folder does not exist")
		}

		fileNames := ctx.Args().Slice()
		if len(fileNames) == 0 {
			return fmt.Errorf("no transfer files available")
		}

		switch network {
		case "tcp":
			tcpAddr, err := net.ResolveTCPAddr(network, addr)
			if err != nil {
				return err
			}

			conTcp, err := net.DialTCP(network, nil, tcpAddr)
			if err != nil {
				return err
			}

			tcpCon := server.NewTcp(conTcp)
			srv := file.NewClient(tcpCon, dir, fileNames)
			_ = srv.SendFile()
		case "udp":
		default:
			return fmt.Errorf("network param err: select tcp | udp")
		}
		return nil
	},
}

func init() {
	registerCommand(clientCmd)
}
