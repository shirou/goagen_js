package main

import (
	"github.com/goadesign/goa"
)

// FilesController implements the files resource.
type FilesController struct {
	*goa.Controller
}

// NewFilesController creates a files controller.
func NewFilesController(service *goa.Service) *FilesController {
	return &FilesController{Controller: service.NewController("FilesController")}
}
