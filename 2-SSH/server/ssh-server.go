package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
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

func handleSession(sess ssh.Session) {
	_, err := io.WriteString(sess, fmt.Sprintf("Hello %s\n", sess.User()))

	if err != nil {
		panic("Write String Error: " + err.Error())
	}

	text, err := bufio.NewReader(sess).ReadString('\n')
    if err != nil  {
      panic("GetLines: " + err.Error())
	}

	term := terminal.NewTerminal(sess, "> ")

    for {
      line, err := term.ReadLine()
      if err != nil {
          break
      }
      log.Println(line)
    }

    io.WriteString(sess, fmt.Sprintf("%s\n", text))
}

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	config := &ssh.Server{
		Addr: addr,
		Handler: handleSession,
		PasswordHandler: passwordHandler,
	}

	fmt.Printf("Opening server on \n" + addr)

	log.Fatal(config.ListenAndServe())
}
