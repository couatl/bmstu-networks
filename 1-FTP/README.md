# Лабораторная работа № 1

* Разработка ftp клиента.
* Разработка ftp сервера.
* Совместная работа ftp сервера и ftp клиента.

## Запуск

Для ftp-client:
```go
export GOPATH=~/go
github.com/secsy/goftp
go run ftp-client.go -user admin -pass 12345 -host localhost
```

Для ftp-server:
```go
export GOPATH=~/go
go get github.com/goftp/file-driver
go get github.com/goftp/server
go run ftp-server.go -root ./root
```
