package main

import (
	"github.com/gliderlabs/ssh"
	"io"
	"log"
	"flag"
	"fmt"
)

var (
  host = flag.String("h", "localhost", "Host")
  port = flag.Int("p", 2216, "Port")
)

func handleSession(s ssh.Session) {
	io.WriteString(s, "Hello world\n")
}

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	config := &ssh.Server{
		Addr: addr,
		Handler: handleSession,
	}

	log.Fatal(config.ListenAndServe())
}
