package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/ywanbing/ft/pkg/file"
	"github.com/ywanbing/ft/pkg/server"
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
			Value: "0.0.0.0:9988",
		},
		&cli.StringFlag{
			Name:    "dir",
			Value:   "./data",
			Aliases: []string{"d"},
			Usage:   "upload dir or save dir",
		},
	},
	Action: func(ctx *cli.Context) error {
		network := ctx.String("network")
		addr := ctx.String("addr")
		dir := ctx.String("dir")

		if !file.PathExists(dir) {
			_ = os.MkdirAll(dir, os.ModePerm)
		}

		switch network {
		case "tcp":
			tcpAddr, err := net.ResolveTCPAddr(network, addr)
			if err != nil {
				return err
			}

			listener, err := net.ListenTCP(network, tcpAddr)
			if err != nil {
				return err
			}

			for {
				acceptTCP, err := listener.AcceptTCP()
				if err != nil {
					return err
				}

				fmt.Println("start a connection:", acceptTCP.RemoteAddr())
				tcpCon := server.NewTcp(acceptTCP)
				srv := file.NewServer(tcpCon, dir)
				_ = srv.Start()
				fmt.Println("end a connection:", acceptTCP.RemoteAddr())
			}
		case "udp":
		default:
			return fmt.Errorf("network param err: select tcp | udp")
		}

		return nil
	},
}
