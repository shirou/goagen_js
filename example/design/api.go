package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("goagen_js", func() {
	Title("Goa gen JS sample")
	Contact(func() {
		Name("shirou")
		Email("shirou.faw@gmail.com")
		URL("https://github.com/shirou/goagen_js")
	})
	Docs(func() {
		Description("goagen js sample")
		URL("https://github.com/shirou/goagen_js")
	})
	Host("localhost:8080")
	Scheme("http")
	Origin("http://localhost:8080", func() {
		Methods("GET", "POST", "PUT", "PATCH", "DELETE")
		MaxAge(600)
		Credentials()
	})

	defineUserTrait()
})
