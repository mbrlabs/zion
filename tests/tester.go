package main

import (
	"hodor"
	"net/http"
	"fmt"
)

type Middleware1 struct {

} 

func (m *Middleware1) Execute(ctx *hodor.Context) {
	fmt.Println("executing middleware 1")
} 

func main() {
	app := hodor.NewHodor()
	mw := new(Middleware1)
	app.MountBefore("", mw)

	app.Get("/test/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "It works!")
	})

	app.Start()
}