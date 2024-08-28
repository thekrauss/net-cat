package main

import (
	"fmt"
	"net-cat/server"
	"os"
	"strings"
)

const Tux = `
Welcome to TCP-Chat!

         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    ` + "`.       | `" + `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     ` + "-'       `" + `--'`

const (
	IP   = "127.0.0.1"
	PORT = "8989"
)

func usage() {
	colorReset := "\033[0m"
	colorBlue := "\033[34m"

	fmt.Println(colorBlue, "exemple pour lancer le server : go run  server", colorReset)
	fmt.Println(colorBlue, "exemple pour lancer client : go run  client", colorReset)
}

func Option() {
	if len(os.Args) != 2 {
		usage()
		os.Exit(2)
	}

	mode := strings.ToLower(os.Args[1])

	if mode == "server" {
		server := server.State(IP, PORT)
		server.Run()
	} else if mode == "client" {
		client := server.State(IP, PORT)
		client.Run()
	} else {
		usage()
		os.Exit(2)
	}
}

func main() {
	fmt.Println(Tux)
	Option()
}
