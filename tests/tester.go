package main

import (
	"hodor"
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

	app.Get("/test/", func(ctx *hodor.Context) {
		fmt.Fprintf(ctx.Writer, "It works!")
	})

	app.Start()
}