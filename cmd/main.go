package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/ywanbing/ft/internal"
)

var help = `

example :
	start server : ./ft -new server -dir ./ -addr 0.0.0.0:8000
	start client : ./ft -new client -dir ./ -addr 127.0.0.1:8000 -file 111.txt
	
note :
	When using client, if dir is not empty,will go down the path of dir to find file`

func main() {
	ne := flag.String("new", "", "new a client or server")
	dir := flag.String("dir", "", "upload dir or download dir")
	addr := flag.String("addr", "", "client con address or server listen address")
	file := flag.String("file", "", "client upload file")
	flag.Parse()

	if *ne == "" || *addr == "" {
		fmt.Println(help)
		os.Exit(-1)
	}

	if *ne == "client" && *file == "" {
		fmt.Println(help)
		os.Exit(-1)
	}

	if *ne == "server" && *dir == "" {
		fmt.Println(help)
		os.Exit(-1)
	}

	if *ne == "server" {
		internal.StartServer(*addr, *dir)
	}

	if *ne == "client" {
		fileName := *file
		if *dir != "" {
			fileName = *dir + "/" + *file
		}
		if err := internal.StartClient(*addr, path.Clean(fileName)); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("send ok !")
		}
	}
}
