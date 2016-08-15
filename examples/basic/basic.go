package main

import (
	"hodor"
	"fmt"
)

type TestMiddleware struct {} 

func (m *TestMiddleware) Execute(ctx *hodor.Context) bool {
	fmt.Printf("Executing %s\n", m.Name())
	return true
} 

func (m *TestMiddleware) Name() string {
	return "Test middleware"
} 

func main() {
	app := hodor.NewHodor()
	app.MountBefore("", new(TestMiddleware))
	app.MountAfter("", new(TestMiddleware))

	app.Get("/test/", func(ctx *hodor.Context) {
		fmt.Fprintf(ctx.Writer, "It works!")
	})

	app.Start()
}