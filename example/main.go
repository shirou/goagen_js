//go:generate goagen bootstrap -d github.com/shirou/goagen_js/example/design

package main

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/shirou/goagen_js/example/app"
)

func main() {
	// Create service
	service := goa.New("goagen_js")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "files" controller
	c := NewFilesController(service)
	app.MountFilesController(service, c)
	// Mount "user" controller
	c2 := NewUserController(service)
	app.MountUserController(service, c2)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}

}
