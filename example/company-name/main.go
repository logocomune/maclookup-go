package main

import (
	"log"

	"github.com/logocomune/maclookup-go"
)

func main() {
	client := maclookup.New()
	r, err := client.CompanyName("000000")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("MAC found in database:", r.Found)
	log.Println("MAC is private (no company name):", r.IsPrivate)
	log.Println("Company name:", r.Company)
	log.Println("Api response in: ", r.RespTime)
	log.Println("Rate limits - remaining request for current time window:", r.RateLimit.Remaining)
	log.Println("Rate limits - next reset", r.RateLimit.Reset)
}
