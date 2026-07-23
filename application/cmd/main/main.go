package main

import (
	"log"

	"github.com/egot3/fathom/cmd"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err == nil {
		log.Printf("running in binary")
	}
	cmd.LoadPath()
	cmd.Execute()
}
