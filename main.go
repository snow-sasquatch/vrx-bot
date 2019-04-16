package main

import (
	"net/http"
	"vrx-bot/providers"
)

func main() {
	c := &http.Client{}
	b := providers.NewBadoinkProvider(c)
	b.Content()
}
