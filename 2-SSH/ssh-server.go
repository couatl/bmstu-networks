package main

import (
	"github.com/gliderlabs/ssh"
	"io"
	"log"
)

var (
  user = flag.String("user", "admin", "Username for login")
  pass = flag.String("pass", "12345", "Password for login")
  host = flag.String("h", "localhost", "Host")
  port = flag.Int("p", 22, "Port")
)

func handleSession(s ssh.Session) {
  io.WriteString(s, "Hello world\n")
}

 func main() {
 		 s := &ssh.Server{
    	Addr:             ":2222",
    	Handler:          sessionHandler,
    	PublicKeyHandler: authHandler,
		 }

     ssh.Handle(handleSession)  

     log.Fatal(ssh.ListenAndServe(":2222", nil))
 }