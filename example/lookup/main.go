package main

import (
	"log"

	"github.com/logocomune/maclookup-go"
)

func main() {
	client := maclookup.New()
	r, err := client.Lookup("000000")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", r)
}
