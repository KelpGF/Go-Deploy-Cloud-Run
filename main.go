package main

import (
	"github.com/KelpGF/Go-Deploy-Cloud-Run/internal/server"
)

func main() {
	server := server.ServerHttp{}

	server.Run()
}
