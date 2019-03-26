package main
import (
"flag"
"fmt"
"io/ioutil"
"github.com/secsy/goftp"
"time"
"os"
)


func main() {
	var (
		user = flag.String("user", "admin", "Username for login")
		pass = flag.String("pass", "12345", "Password for login")
		host = flag.String("host", "localhost", "FTP server to connect")
	)
	flag.Parse()

	config := goftp.Config {
	    User:               *user,
	    Password:           *pass,
	    ConnectionsPerHost: 10,
	    Timeout:            10 * time.Second,
	    Logger:             os.Stderr,
	}

	ftp_server := *host
	client, err := goftp.DialConfig(config, ftp_server)

	if err != nil {
	    panic(err)
	}

	rawConn, _ := client.OpenRawConn()

	fmt.Printf("Connected to: %s\n", ftp_server)
	fmt.Printf("pass: %s, user: %s, host: %s\n", *pass, *user, *host)

	var s string
	var params string
	var params_second string

	for {
		fmt.Printf("Enter help to see available commands. \n")
		fmt.Printf("Enter your command: ")
		fmt.Scanf("%s %s %s", &s, &params, &params_second)

		if s == "ls" {
			if params == "" {
				dcGetter, _ := rawConn.PrepareDataConn()
				rawConn.SendCommand("LIST")
				dc, _ := dcGetter()
				data, _ := ioutil.ReadAll(dc)
				fmt.Printf("LIST response: %s\n", data)
				dc.Close()
				code, msg, _ := rawConn.ReadResponse()
				fmt.Printf("Final response: %d-%s\n", code, msg)
			} else {
				// получение содержимого директории на ftp сервере с помощью go ftp клиента.
				res, err := client.ReadDir(params)
				if err != nil {
					continue
				}
				
			    for _, f := range res {
			            fmt.Println(f.Name())
			    }
			}
		} else if s == "exit" {
			os.Exit(2)
		} else if s == "mkdir" {
			// создание директории go ftp клиентом на ftp сервере
			fmt.Printf("*** Creating directory %s \n", params)
			res, err := client.Mkdir(params)
			if err != nil {
				continue
			}
			fmt.Printf("- Response: %s \n", res)
		} else if s == "copy" {
			// загрузка файла go ftp клиентом на ftp сервер
			fmt.Printf("*** Copying file %s to %s \n", params, params_second)

			file, err := os.Open(params)
			if err != nil {
			    continue
			}

			err = client.Store(params_second, file)
			if err != nil {
			    continue
			}
			fmt.Printf("- Successfull \n")
		} else if s == "rm" {
			err = client.Delete(params)
			if err != nil {
			    continue
			}
			fmt.Printf("- Successfull \n")
		} else if s == "dwnld" {
			fmt.Printf("*** Downloading file %s to %s \n", params, params_second)

			file, err := os.Create(params_second)
			if err != nil {
			    continue
			}

			err = client.Retrieve(params_second, file)
			if err != nil {
			    continue
			}
			fmt.Printf("- Successfull \n")
		} else {
			fmt.Printf("Unexpected command. Type exit to exit \n")
		}
	}
	
}