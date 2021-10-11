package main

import (
	"net/http"

	"github.com/mkulke/go-openapi-playground/api"
)

func main() {
  api.SetupHandler()
  http.ListenAndServe(":8080", nil)
}
