package main

import (
	"log"
	"syscall"

	platform "blaucorp.com/fiscal-engine/internal/service"
)

const name = "FiscalEngine"

func main() {

	if err := platform.Engage(); err != nil {
		log.Fatalf("%s service struggles with (%v)", name, err)
	}

	syscall.Exit(0)
}
