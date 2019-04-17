package main

import (
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"flag"
	"fmt"
	"bufio"
)

var (
  host = flag.String("h", "localhost", "Host")
  port = flag.Int("p", 2216, "Port")
  pass = flag.String("pass", "12345", "Password for login")
  user = flag.String("user", "admin", "Username for login")
)

func passwordHandler(ctx ssh.Context, password string) bool {
	return ctx.User() == *user && *pass == password
}

func handleSession(s ssh.Session) {
	io.WriteString(s, fmt.Sprintf("Hello %s\n", s.User()))

	text, err:= bufio.NewReader(s).ReadString('\n')
    if err != nil  {
        panic("GetLines: " + err.Error())
	}
	
	term := terminal.NewTerminal(sess, "> ")
    for {
        line, err := term.ReadLine()
        if err != nil {
            break
        }
        response := router(line)
        log.Println(line)
        if response != "" {
            term.Write(append([]byte(response), '\n'))
        }
    }

    io.WriteString(s, fmt.Sprintf("%s\n", text))   
}

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	config := &ssh.Server{
		Addr: addr,
		Handler: handleSession,
		PasswordHandler: passwordHandler,
	}

	log.Fatal(config.ListenAndServe())
}
