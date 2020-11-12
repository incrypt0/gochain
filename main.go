package main

import (
	"log"
	"os"

	"github.com/incrypt0/gochain/blockchain"
	"github.com/incrypt0/gochain/handler"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(os.Getenv("PORT"))

	chain := blockchain.New()

	handler.New(chain)
}
