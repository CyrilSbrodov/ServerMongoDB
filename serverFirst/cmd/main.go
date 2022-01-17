package main

import (
	"server/serverFirst/cmd/server"
)


func main() {

	server := server.NewApp()
	server.Run()

}

