package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	user = flag.String("user", "admin", "Username for login")
	pass = flag.String("pass", "12345", "Password for login")
	host = flag.String("h", "localhost", "Host")
	port = flag.Int("p", 22, "Port")
)

func executeCmd(cmd, addr string, config *ssh.ClientConfig) string {
	connection, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		panic(err)
	}

	session, err := connection.NewSession()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	fmt.Printf("Connected to: %s@%s:%d\n", *user, *host, *port)

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(cmd)

	return addr + ": " + stdoutBuf.String()
}

func main() {
	flag.Parse()

	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			ssh.Password(*pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var cmd string
	var params string
	addr := fmt.Sprintf("%s:%d", *host, *port)

	for {
		fmt.Printf("exit for exit \n")
		fmt.Printf("Enter your command: ")
		fmt.Scanf("%s %s", &cmd, &params)

		if cmd == "exit" {
			fmt.Printf("Exited.")
			os.Exit(2)
		} else {
			command := fmt.Sprintf("%s %s\n", cmd, params)
			fmt.Printf(command)

			result := executeCmd(command, addr, config)
			fmt.Print(result)
		}
	}
}
