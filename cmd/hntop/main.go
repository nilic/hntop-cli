package main

import (
	"log"
	"os"
)

func main() {
	app := newApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
