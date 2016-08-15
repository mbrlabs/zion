package main

import (
	"hodor"
	"net/http"
	"fmt"
)

func main() {
	hodor := hodor.NewHodor()
	hodor.Use(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("running middleware 1")
	})
	hodor.Use(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("running middleware 2")
	})
	hodor.Use(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("running middleware 3")
	})

	hodor.Get("/test/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "It works!")
	})

	hodor.Start()
}