package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/logocomune/maclookup-go"
	"golang.org/x/time/rate"
)

const (
	maxReqSec = 2
)

func main() {
	limiter := rate.NewLimiter(rate.Every(time.Second/maxReqSec), 1)

	var wg sync.WaitGroup

	client := maclookup.New()
	wg.Add(maxReqSec * 2)

	for i := 0; i < maxReqSec*2; i++ {
		go func() {
			defer wg.Done()

			err := limiter.Wait(context.Background())
			if err != nil {
				fmt.Printf("rate limit error: %v", err)
			}

			r, err := client.CompanyName("00:00:00")
			if err != nil {
				var e *maclookup.RateLimitsExceeded
				if errors.As(err, &e) {
					log.Println("Rate limits exceeded", err.Error())
				}
			} else {
				log.Println(r.Company, "- (response in", r.RespTime, ")")
			}
		}()
	}
	wg.Wait()
}
