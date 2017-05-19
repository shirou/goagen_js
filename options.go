package goagen_js

import "github.com/goadesign/goa/design"

//Option a generator option definition
type Option func(*Generator)

//API The API definition
func API(API *design.APIDefinition) Option {
	return func(g *Generator) {
		g.API = API
	}
}

//OutDir Path to output directory
func OutDir(outDir string) Option {
	return func(g *Generator) {
		g.OutDir = outDir
	}
}

//Scheme Scheme used by JavaScript client
func Scheme(scheme string) Option {
	return func(g *Generator) {
		g.Scheme = scheme
	}
}

//Host addressed by JavaScript client
func Host(host string) Option {
	return func(g *Generator) {
		g.Host = host
	}
}
