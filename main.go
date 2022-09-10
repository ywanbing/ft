package main

import (
	"log"
	"os"

	"github.com/ywanbing/ft/cmd"
)

func main() {
	if err := cmd.NewApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
