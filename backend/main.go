package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"
)

// @title Airi-Go
// @version 1.0
// @description Airi Self hosted, you owned Grok Companion, a container of souls of waifu, cyber livings to bring them into our worlds.
// @contact.name kiosk
// @contact.url http://www.swagger.io/support
// @contact.email kiosk007@gmail.com
func main() {
	rand.Seed(time.Now().UnixNano())
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	command := NewAiriGoCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
