package main

import (
	"server/internal/server"
	"time"
)

func main() {

	time.Sleep(3 * time.Second)
	server.Run()
}
