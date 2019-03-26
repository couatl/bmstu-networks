# Лабораторная работа № 1

* Разработка ftp клиента.
* Разработка ftp сервера.
* Совместная работа ftp сервера и ftp клиента.

## Запуск

Для ftp-client:
```go
export GOPATH=~/go
github.com/secsy/goftp
go run main.go
```

Для ftp-server:
```go
export GOPATH=~/go
go get github.com/goftp/file-driver
go get github.com/goftp/server
go run main.go
```
