package goagen_js

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
	"github.com/goadesign/goa/goagen/utils"
)

//NewGenerator returns an initialized instance of a JavaScript Client Generator
func NewGenerator(options ...Option) *Generator {
	g := &Generator{}

	for _, option := range options {
		option(g)
	}

	return g
}

// Generator is the application code generator.
type Generator struct {
	API      *design.APIDefinition // The API definition
	OutDir   string                // Destination directory
	Scheme   string                // Scheme used by JavaScript client
	Host     string                // Host addressed by JavaScript client
	Target   string                // Target JS (es2015, flowtype, tc)
	genfiles []string              // Generated files
}

// Generate is the generator entry point called by the meta generator.
func Generate() (files []string, err error) {
	var (
		outDir, ver  string
		scheme, host string
	)

	set := flag.NewFlagSet("client", flag.PanicOnError)
	set.StringVar(&outDir, "out", "", "")
	set.String("design", "", "")
	set.StringVar(&scheme, "scheme", "", "")
	set.StringVar(&host, "host", "", "")
	set.StringVar(&ver, "version", "", "")
	set.Parse(os.Args[1:])

	// First check compatibility
	if err := codegen.CheckVersion(ver); err != nil {
		return nil, err
	}

	// Now proceed
	g := &Generator{
		OutDir: outDir,
		Scheme: scheme,
		Host:   host,
		API:    design.Design,
	}

	return g.Generate()
}

// Generate produces the skeleton main.
func (g *Generator) Generate() ([]string, error) {
	var err error
	if g.API == nil {
		return nil, fmt.Errorf("missing API definition, make sure design is properly initialized")
	}

	go utils.Catch(nil, func() { g.Cleanup() })

	defer func() {
		if err != nil {
			g.Cleanup()
		}
	}()

	if g.Scheme == "" && len(g.API.Schemes) > 0 {
		g.Scheme = g.API.Schemes[0]
	}
	if g.Scheme == "" {
		g.Scheme = "http"
	}
	if g.Host == "" {
		g.Host = g.API.Host
	}
	if g.Host == "" {
		return nil, fmt.Errorf("missing host value, set it with --host")
	}

	g.OutDir = filepath.Join(g.OutDir, "js")
	if err := os.MkdirAll(g.OutDir, 0755); err != nil {
		return nil, err
	}

	ps, err := parseActions(g)
	if err != nil {
		return nil, err
	}

	// Generate api_request.js
	if err := g.generateRequestJS(filepath.Join(g.OutDir, "api_request.js"), ps); err != nil {
		return nil, err
	}

	// Generate validate.js
	if err := g.generateValidateJS(filepath.Join(g.OutDir, "api_validator.js"), ps); err != nil {
		return nil, err
	}

	return g.genfiles, nil
}

func (g *Generator) generateRequestJS(jsFile string, params []ParamsDefinition) error {
	// clean exist file
	if err := ensureDelete(jsFile); err != nil {
		return err
	}

	file, err := codegen.SourceFileFor(jsFile)
	if err != nil {
		return err
	}
	g.genfiles = append(g.genfiles, jsFile)

	data := map[string]interface{}{
		"API":    g.API,
		"Host":   g.Host,
		"Scheme": g.Scheme,
	}
	if err = file.ExecuteTemplate("header", jsHeaderT, nil, data); err != nil {
		return err
	}

	for _, p := range params {
		data := map[string]interface{}{
			"Action":           p.Action,
			"FuncName":         p.FuncName(),
			"Args":             p.FuncArgs(g.Target),
			"UrlArgs":          p.UrlArgs(),
			"PathParams":       p.Path,
			"QueryParams":      p.Query,
			"Comments":         p.Comments(p.Action),
			"Request":          p.Request(),
			"ValidateRequired": p.ValidateRequired(),
		}
		funcs := template.FuncMap{}
		if err = file.ExecuteTemplate("jsFuncs", jsFuncsT, funcs, data); err != nil {
			return err
		}
	}
	if err = file.ExecuteTemplate("module", jsModuleT, nil, data); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateValidateJS(jsFile string, params []ParamsDefinition) error {
	// clean exist file
	if err := ensureDelete(jsFile); err != nil {
		return err
	}
	file, err := codegen.SourceFileFor(jsFile)
	if err != nil {
		return err
	}
	g.genfiles = append(g.genfiles, jsFile)

	data := map[string]interface{}{
		"API": g.API,
	}
	if err := file.ExecuteTemplate("header", validatorHeaderT, nil, data); err != nil {
		return err
	}

	for _, p := range params {
		json, err := p.Validator.JSONify()
		if err != nil {
			return err
		}
		data := map[string]interface{}{
			"Name":       p.FuncName(),
			"Constraint": json,
		}
		funcs := template.FuncMap{}
		if err := file.ExecuteTemplate("valid", validatorT, funcs, data); err != nil {
			return err
		}
	}
	if err := file.ExecuteTemplate("module", validatorModuleT, nil, data); err != nil {
		return err
	}

	return nil
}

// Cleanup removes all the files generated by this generator during the last invokation of Generate.
func (g *Generator) Cleanup() {
	for _, f := range g.genfiles {
		os.Remove(f)
	}
	g.genfiles = nil
}

func ensureDelete(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

func getActions(g *Generator) map[string][]*design.ActionDefinition {
	actions := make(map[string][]*design.ActionDefinition)
	g.API.IterateResources(func(res *design.ResourceDefinition) error {
		return res.IterateActions(func(action *design.ActionDefinition) error {
			if as, ok := actions[action.Name]; ok {
				actions[action.Name] = append(as, action)
			} else {
				actions[action.Name] = []*design.ActionDefinition{action}
			}
			return nil
		})
	})
	return actions
}