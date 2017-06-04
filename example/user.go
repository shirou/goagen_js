package main

import (
	"github.com/goadesign/goa"
	"github.com/shirou/goagen_js/example/app"
)

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service) *UserController {
	return &UserController{Controller: service.NewController("UserController")}
}

// Create runs the create action.
func (c *UserController) Create(ctx *app.CreateUserContext) error {
	// UserController_Create: start_implement

	// Put your logic here

	// UserController_Create: end_implement
	return nil
}

// Get runs the get action.
func (c *UserController) Get(ctx *app.GetUserContext) error {
	// UserController_Get: start_implement

	// Put your logic here

	// UserController_Get: end_implement
	res := &app.User{}
	return ctx.OK(res)
}

// List runs the list action.
func (c *UserController) List(ctx *app.ListUserContext) error {
	// UserController_List: start_implement

	// Put your logic here

	// UserController_List: end_implement
	res := app.UserCollection{}
	return ctx.OK(res)
}
