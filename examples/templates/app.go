/*package main

import (
	"github.com/mbrlabs/hodor"
)

func main() {
	app := hodor.NewHodor()
	app.ServeStaticFiles("/static/*", "static/")
	app.Templates("html/")

	// index
	app.Get("/", func(ctx *hodor.Context) {
		ctx.Render("index", "data: index")
	})

	// home
	app.Get("/home/", func(ctx *hodor.Context) {
		ctx.Render("home", "data: home")
	})

	app.Start()
}
*/