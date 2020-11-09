package main

import (
	"log"
	"os"

	"github.com/incrypt0/gochain/handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	h := handler.New()
	go h.GenesisBlock()
	h.Register(e)

	log.Println(os.Getenv("PORT"))

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
