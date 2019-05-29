package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type Mail struct {
	toIds    []string
	subject  string
	body     string
}

type Connection struct {
	host string `json:"host"`
	port string `json:"port"`
	login string `json:"login"`
}

func (s *Connection) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage(login string) string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", login)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return decrypt(data, passphrase)
}

func getPasswd() string {
	var key string
	fmt.Printf("Enter key to your passphrase: ")
	fmt.Scanf("%s", &key)

	cipherpass, _ := ioutil.ReadFile(*passwdFile)
	passwd := string(decrypt(cipherpass, key))
	return passwd
}

var (
	paramsFile = flag.String("params", "params", "File of connection config")
	passwdFile = flag.String("pass", "passwd", "File with encrypted password")
)

func main() {
	flag.Parse()

	mail := Mail{}
	connection := Connection{}
	data, err := ioutil.ReadFile(*paramsFile)
	if err != nil {
		panic(err.Error())
	}

	params := strings.Split(string(data), "\n")
	var param []string
	for _, paramStr := range params {
		param = strings.Split(paramStr, ":")
		log.Println(param)
		if param[0] == "host" {
			connection.host = param[1]
		} else if param[0] == "port" {
			connection.port = param[1]
		} else if param[0] == "login" {
			connection.login = param[1]
		}
	}

	for {
		var cmd string

		fmt.Printf("exit for exit \n")
		fmt.Printf("Enter anything to start making a message: ")
		fmt.Scanf("%s", &cmd)

		if cmd == "exit" {
			fmt.Printf("Exited.")
			os.Exit(2)
		} else if cmd == "pass" {
			var password string
			var key string
			fmt.Printf("Enter key to your passphrase: ")
			fmt.Scanf("%s", &key)

			fmt.Printf("Enter your passphrase: ")
			fmt.Scanf("%s", &password)

			ciphertext := encrypt([]byte(password), key)
			ioutil.WriteFile(*passwdFile, ciphertext, 0644)
		} else {
			var to string

			fmt.Printf("Enter To separated by comma: ")
			fmt.Scanf("%s", &to)

			fmt.Printf("Enter Subject: ")
			fmt.Scanf("%s", &mail.subject)

			fmt.Printf("Enter Body: ")
			fmt.Scanf("%s", &mail.body)

			mail.toIds = strings.Split(to, ",")
			messageBody := mail.BuildMessage(connection.login)

			log.Println(connection.ServerName(), messageBody)

			auth := smtp.PlainAuth("", connection.login, getPasswd(), connection.host)

			tlsConfig := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         connection.host,
			}

			conn, err := tls.Dial("tcp", connection.ServerName(), tlsConfig)
			if err != nil {
				log.Panic(err)
			}

			client, err := smtp.NewClient(conn, connection.host)
			if err != nil {
				log.Panic(err)
			}

			if err = client.Auth(auth); err != nil {
				log.Panic(err)
			}

			if err = client.Mail(connection.login); err != nil {
				log.Panic(err)
			}
			for _, k := range mail.toIds {
				if err = client.Rcpt(k); err != nil {
					log.Panic(err)
				}
			}

			w, err := client.Data()
			if err != nil {
				log.Panic(err)
			}

			_, err = w.Write([]byte(messageBody))
			if err != nil {
				log.Panic(err)
			}

			err = w.Close()
			if err != nil {
				log.Panic(err)
			}

			client.Quit()

			log.Println("Mail sent successfully")
		}
	}
}
