package main

import (
	"net/http"
	"vrx-bot/badoink"
)

func main() {
	c := &http.Client{}
	b := badoink.NewProvider(c)
	b.Content()
}
