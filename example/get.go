package main

import (
	"github.com/goadesign/goa"
	"github.com/shirou/goagen_js/example/app"
)

// GetController implements the get resource.
type GetController struct {
	*goa.Controller
}

// NewGetController creates a get controller.
func NewGetController(service *goa.Service) *GetController {
	return &GetController{Controller: service.NewController("GetController")}
}

// GetInt runs the get_int action.
func (c *GetController) GetInt(ctx *app.GetIntGetContext) error {
	// GetController_GetInt: start_implement

	// Put your logic here

	// GetController_GetInt: end_implement
	return nil
}

// PathParams runs the path_params action.
func (c *GetController) PathParams(ctx *app.PathParamsGetContext) error {
	// GetController_PathParams: start_implement

	// Put your logic here

	// GetController_PathParams: end_implement
	res := &app.Inttest{}
	return ctx.OK(res)
}

// Without runs the without action.
func (c *GetController) Without(ctx *app.WithoutGetContext) error {
	// GetController_Without: start_implement

	// Put your logic here

	// GetController_Without: end_implement
	res := &app.Inttest{}
	return ctx.OK(res)
}
