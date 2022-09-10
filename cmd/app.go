package cmd

import (
	"github.com/urfave/cli/v2"
)

var commands []*cli.Command

// 提供其他命令注册
func registerCommand(cmd *cli.Command) {
	commands = append(commands, cmd)
}

// NewApp 创建一个 cli APP，并组装所有的命令。
func NewApp() *cli.App {
	app := &cli.App{
		Name:  "ft",
		Usage: "big file transfer, support various network protocols",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Value:   "./data",
				Aliases: []string{"d"},
				Usage:   "upload dir or save dir",
			},
		},
		Commands:             commands,
		EnableBashCompletion: true,
	}

	return app
}
