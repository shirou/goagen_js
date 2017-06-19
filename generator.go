package goagen_js

import (
	"bufio"
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
		outDir, ver    string
		scheme, host   string
		target, genOut string
	)

	set := flag.NewFlagSet("client", flag.PanicOnError)
	set.StringVar(&outDir, "out", "", "")
	set.String("design", "", "")
	set.StringVar(&scheme, "scheme", "", "")
	set.StringVar(&host, "host", "", "")
	set.StringVar(&ver, "version", "", "")
	set.StringVar(&target, "target", "js", "target language(js, flow, type)")
	set.StringVar(&genOut, "genout", "js", "output directory")
	set.Parse(os.Args[1:])

	// First check compatibility
	if err := codegen.CheckVersion(ver); err != nil {
		return nil, err
	}

	// Now proceed
	g := &Generator{
		OutDir: filepath.Join(outDir, genOut),
		Scheme: scheme,
		Host:   host,
		API:    design.Design,
		Target: target,
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

	if err := os.MkdirAll(g.OutDir, 0755); err != nil {
		return nil, err
	}

	ps, err := parseActions(g)
	if err != nil {
		return nil, err
	}
	requestJS := "api_request.js"
	validatorJS := "api_validator.js"
	definitionJS := ""

	switch g.Target {
	case TargetFlow:
		definitionJS = "api.d.js"
	case TargetTS:
		requestJS = "api_request.ts"
		validatorJS = "api_validator.ts"
		definitionJS = "api.d.ts"
	}

	// Generate api_request.js
	if err := g.generateRequestJS(filepath.Join(g.OutDir, requestJS), ps); err != nil {
		return nil, err
	}

	// Generate validate.js
	if err := g.generateValidateJS(filepath.Join(g.OutDir, validatorJS), ps); err != nil {
		return nil, err
	}

	// Generate definition
	if definitionJS != "" {
		if err := g.generateDefinition(filepath.Join(g.OutDir, definitionJS), ps); err != nil {
			return nil, err
		}
	}

	return g.genfiles, nil
}

func (g *Generator) generateRequestJS(jsFile string, params []ParamsDefinition) error {
	file, err := openFile(jsFile)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := bufio.NewWriter(file)
	g.genfiles = append(g.genfiles, jsFile)

	tmpl, err := newTemplate(
		jsHeaderT,
		jsFuncsT,
		jsModuleT,
	)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"API":    g.API,
		"Host":   g.Host,
		"Scheme": g.Scheme,
		"Target": g.Target,
	}
	// write header
	if err := tmpl.ExecuteTemplate(buf, "js_header", data); err != nil {
		return err
	}
	for _, p := range params {
		data := map[string]interface{}{
			"Action":           p.Action,
			"FuncName":         p.FuncName(),
			"Args":             p.FuncArgs(g.Target),
			"FuncRet":          p.FuncRet(g.Target),
			"UrlArgs":          p.UrlArgs(),
			"PathParams":       p.Path,
			"QueryParams":      p.Query,
			"Comments":         p.Comments(p.Action),
			"Request":          p.Request(),
			"ValidateRequired": p.ValidateRequired(),
		}
		if err := tmpl.ExecuteTemplate(buf, "js_funcs", data); err != nil {
			return err
		}
	}
	if err := tmpl.ExecuteTemplate(buf, "js_module", data); err != nil {
		return err
	}
	return buf.Flush()
}

func (g *Generator) generateValidateJS(jsFile string, params []ParamsDefinition) error {
	file, err := openFile(jsFile)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := bufio.NewWriter(file)
	g.genfiles = append(g.genfiles, jsFile)

	data := map[string]interface{}{
		"API":    g.API,
		"Target": g.Target,
	}

	tmpl, err := newTemplate(
		validatorHeaderT,
		validatorT,
		validatorModuleT,
	)
	if err != nil {
		return err
	}

	if err := tmpl.ExecuteTemplate(buf, "validator_header", data); err != nil {
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
		if err := tmpl.ExecuteTemplate(buf, "validator_definition", data); err != nil {
			return err
		}
	}
	if err := tmpl.ExecuteTemplate(buf, "validator_module", data); err != nil {
		return err
	}

	return buf.Flush()
}

func (g *Generator) generateDefinition(jsFile string, params []ParamsDefinition) error {
	file, err := openFile(jsFile)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := bufio.NewWriter(file)
	g.genfiles = append(g.genfiles, jsFile)

	data := map[string]interface{}{
		"API":    g.API,
		"Target": g.Target,
	}

	var tmpl *template.Template
	switch g.Target {
	case TargetFlow:
		tmpl, err = newTemplate(
			definitionHeaderFlow,
			definitionFlow,
		)
	case TargetTS:
		tmpl, err = newTemplate(
			definitionHeaderType,
			definitionType,
		)
	default:
		return fmt.Errorf("unknown type language: %s", g.Target)
	}
	if err != nil {
		return err
	}

	if err := tmpl.ExecuteTemplate(buf, "definition_header", data); err != nil {
		return err
	}
	for _, p := range params {
		if len(p.Query) > 0 {
			data := map[string]interface{}{
				"Name":       p.FuncName() + "Payload",
				"Definition": p.PayloadDefinition(g.Target),
			}
			if err := tmpl.ExecuteTemplate(buf, "definition", data); err != nil {
				return err
			}
		}
		if p.Response != nil {
			// generate Response MediaType
			data := map[string]interface{}{
				"Name":       p.Response.IdentifierName,
				"Definition": p.ResponseDefinition(g.Target),
			}
			if err := tmpl.ExecuteTemplate(buf, "definition", data); err != nil {
				return err
			}
		}
	}

	return buf.Flush()
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

func openFile(path string) (*os.File, error) {
	// clean exist file
	if err := ensureDelete(path); err != nil {
		return nil, err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	file, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return file, nil
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
