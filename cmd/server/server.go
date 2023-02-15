package main

import (
	"encoding/json"
	"github.com/bluntenpassant/word-of-wisdom-tcp/manager"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	rawQuotes, err := os.ReadFile("./static/quotes.json")
	if err != nil {
		panic(err)
	}

	var quotes []string

	err = json.Unmarshal(rawQuotes, &quotes)
	if err != nil {
		panic(err)
	}

	connManager := manager.NewManager(quotes)
	err = connManager.Run("0.0.0.0:8081")
	if err != nil {
		panic(err)
	}
}
