package goagen_js

import (
	"fmt"
	"sort"
	"strings"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
)

type ParamsDefinition struct {
	Action    *design.ActionDefinition
	Base      string
	Name      string
	Path      Params // sorted by goa order
	Query     Params // sorted by alphabetical
	Validator Validator
}

func parseActions(g *Generator) ([]ParamsDefinition, error) {
	ret := make([]ParamsDefinition, 0)
	actions := getActions(g)

	keys := []string{}
	for n := range actions {
		keys = append(keys, n)
	}
	sort.Strings(keys)
	for _, n := range keys {
		for _, a := range actions[n] {
			p, err := parseAction(a, g.Target)
			if err != nil {
				return nil, err
			}
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func parseAction(action *design.ActionDefinition, target string) (ParamsDefinition, error) {
	ret := ParamsDefinition{
		Action:    action,
		Name:      codegen.Goify(action.Name, true),
		Path:      make(Params, 0),
		Query:     make(Params, 0),
		Validator: newValidator(codegen.Goify(action.Name, true)),
	}

	if action.PathParams() != nil {
		m := action.PathParams().Type.ToObject()
		for a, l := range m {
			ret.Path = append(ret.Path, newParam(action, target, a, l))
			kind := convertTypeString(l.Type.Kind(), target)
			ret.Validator.constraint[a] = parseConstraint(kind, l.Validation)
		}
	}

	if action.QueryParams != nil {
		m := action.QueryParams.Type.ToObject()
		pv := make(map[string]Constraint)
		for a, l := range m {
			kind := convertTypeString(l.Type.Kind(), target)
			ret.Query = append(ret.Query, newParam(action, target, a, l))
			pv[a] = parseConstraint(kind, l.Validation)
		}
		if len(pv) > 0 {
			ret.Validator.constraint["payload"] = pv
		}
	}

	return ret, nil
}

func (p ParamsDefinition) Comments(action *design.ActionDefinition) []string {
	c := make([]string, 0)

	c = append(c, fmt.Sprintf("%s", p.Name))
	if p.Action.Description != "" {
		c = append(c, p.Action.Description)
	}
	c = append(c, "")
	for _, path := range p.Path {
		c = append(c, fmt.Sprintf("%s(%s): %s", path.Name, path.Kind, path.Description))
	}

	// TODO: why strings.Join(c, "\n"+`// `) expaneded to "../../../" ?

	return c
}

func (p ParamsDefinition) ValidateRequired() bool {
	if len(p.Path) > 0 || len(p.Query) > 0 {
		return true
	}
	return false
}

func (p ParamsDefinition) FuncName() string {
	var v string
	switch p.Action.Routes[0].Verb {
	case "GET":
		v = "Get"
	case "PUT":
		v = "Put"
	case "POST":
		v = "Post"
	case "DELETE":
		v = "Delete"
	case "OPTIONS":
		v = "Options"
	case "HEAD":
		v = "Head"
	case "PATCH":
		v = "Patch"
	default:
		v = "Get"
	}

	return p.Name + v
}

func (p ParamsDefinition) FuncArgs(target string) string {
	ret := make([]string, len(p.Path))
	for i, p := range p.Path {
		ret[i] = p.Name
	}
	if len(p.Query) > 0 {
		ret = append(ret, "payload")
	}

	return strings.Join(ret, ", ")
}

func (p ParamsDefinition) UrlArgs() string {
	path := p.Action.Routes[0].FullPath()
	for _, o := range p.Action.Routes[0].Params() {
		d := fmt.Sprintf("${%s}", codegen.Goify(o, false))
		path = strings.Replace(path, ":"+o, d, 1)
	}

	return `"` + path + `"`
}

func (p ParamsDefinition) Request() string {
	verb := strings.ToLower(p.Action.Routes[0].Verb)

	args := []string{"url"}
	if len(p.Query) > 0 {
		args = append(args, "payload")
	}

	return fmt.Sprintf("%s(%s)", verb, strings.Join(args, ", "))
}

func convertTypeString(t design.Kind, target string) string {
	switch t {
	case design.BooleanKind:
		return "bool"
	case design.IntegerKind:
		return "number"
	case design.NumberKind:
		return "number"
	case design.StringKind:
		return "string"
	case design.DateTimeKind:
		return "datetime"
	case design.UUIDKind:
		return "string"
	case design.AnyKind:
		return "any"
	case design.ArrayKind:
		return "array"
	case design.ObjectKind:
		return "object"
	case design.HashKind:
		return "hash"
	case design.UserTypeKind:
		return "any"
	case design.MediaTypeKind:
		return "any"
	}
	return "any"
}

type Param struct {
	original    *design.AttributeDefinition
	Name        string // CamelCase name
	Kind        string // kind such as bool, int, ...
	Description string
}

type Params []Param

func newParam(action *design.ActionDefinition, target string, a string, l *design.AttributeDefinition) Param {
	return Param{
		original:    l,
		Name:        codegen.Goify(a, false),
		Kind:        convertTypeString(l.Type.Kind(), target),
		Description: l.Description,
	}
}