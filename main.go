package main

import (
	"log"
	"time"
	// "fmt"
	// "net/http"
	"github.com/gagliardetto/solana-go"
	"go.blockdaemon.com/pyth"
	"github.com/shopspring/decimal"

)

func main() {

	// Connect to Pyth on Solana devnet.
	client := pyth.NewClient(pyth.Devnet, "", "ws://api.devnet.solana.com")

	// Open new event stream.
	stream := client.StreamPriceAccounts()
	handler := pyth.NewPriceEventHandler(stream)
	
	// Subscribe to price account changes.
	priceKey := solana.MustPublicKeyFromBase58("J83w4HKfqxwcq3BEMMkPFSppX3gqekLyLJBexebFVkix")
	prev_price, _ := decimal.NewFromString("30.40")
	zero_price, _ := decimal.NewFromString("0.0")
	prev_price_diff := zero_price 
	percent, _ := decimal.NewFromString("100")
	call_percent, _:= decimal.NewFromString("-30")
	
	
	handler.OnPriceChange(priceKey, func(info pyth.PriceUpdate) {
		price, conf, ok := info.Current()
		if ok {


			// is triggered every 6hrs
			time.AfterFunc(6 * time.Hour, func() {
				log.Printf("Timer is called")
				prev_price_diff = zero_price
				prev_price = price
			})

			price_diff := price.Sub(prev_price).Add(prev_price_diff)
			log.Printf("Price diff: %s", price_diff)
			price_change_percent := (price_diff.Div(price).Mul(percent))
			log.Printf("Price change: $%s ± $%s Change: %s Percent: %s% ", price, conf, price_diff, price_change_percent)
			prev_price = price	
			prev_price_diff = price_diff

			if price_change_percent.LessThanOrEqual(call_percent) {
				log.Printf("Price Change Alert! Price dropped below 30% in 6hrs")
			}

			// priceCheck := func (w ResponseWriter, r *Request)
			// 	{fmt.Fprintf(w, price)}

			// handleForPrice := http.HandlerFunc(priceCheck)

			// http.Handle("/metrics", handleForPrice )
		}
	})
	
    // http.ListenAndServe(":2112", nil)
	// Close stream after a while.
	for true {}
	
	// <-time.After(10 * time.Second)
	stream.Close()
	
	
}