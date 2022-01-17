package main

import (
	server2 "server/serverSecond/cmd/server"
)


func main() {

	server := server2.NewApp()
	server.Run()

}

