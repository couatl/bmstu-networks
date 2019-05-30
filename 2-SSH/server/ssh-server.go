package main

import (
	"flag"
	"fmt"
	"github.com/gliderlabs/ssh"
	"log"
	"net"
	"os/exec"
)

var (
	host = flag.String("h", "", "Host")
	port = flag.Int("p", 2216, "Port")
	pass = flag.String("pass", "12345", "Password for login")
	user = flag.String("user", "admin", "Username for login")
)

func passwordHandler(ctx ssh.Context, password string) bool {
	return ctx.User() == *user && *pass == password
}

func greetingHandler(conn net.Conn) net.Conn {
	conn.Write([] byte("Welcome to server\n"))

	return conn
}

func handleSession(sess ssh.Session) {
	if len(sess.Command()) > 0 {
		cmd := exec.Command(sess.Command()[0], sess.Command()[1:]...)
		stdoutStderr, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		sess.Write(stdoutStderr)
	}
}

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	server := &ssh.Server{
		Addr: addr,
		Handler: handleSession,
		ConnCallback: greetingHandler,
	}

	server.SetOption(ssh.PasswordAuth(passwordHandler))

	fmt.Printf("Opening server on \n" + addr)

	log.Fatal(server.ListenAndServe())
}