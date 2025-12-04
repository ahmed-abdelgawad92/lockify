package main

import (
	"log"

	"github.com/ahmed-abdelgawad92/lockify/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
