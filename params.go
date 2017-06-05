package goagen_js

import (
	"fmt"
	"sort"
	"strings"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
)

const (
	TargetJS   = "js"
	TargetFlow = "flow"
	TargetTS   = "type"
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

	name := codegen.Goify(action.Parent.Name, true) + codegen.Goify(action.Name, true)

	ret := ParamsDefinition{
		Action:    action,
		Name:      name,
		Path:      make(Params, 0),
		Query:     make(Params, 0),
		Validator: newValidator(name),
	}

	if action.PathParams() != nil {
		m := action.PathParams().Type.ToObject()
		for a, att := range m {
			ret.Path = append(ret.Path, newParam(action, target, a, att))
			kind := convertTypeString(att.Type.Kind(), target)
			ret.Validator.constraint[a] = parseConstraint(kind, att.Validation, att.IsRequired(a))
		}
	}

	if action.QueryParams != nil {
		m := action.QueryParams.Type.ToObject()
		pv := make(map[string]Constraint)
		for a, att := range m {
			kind := convertTypeString(att.Type.Kind(), target)
			ret.Query = append(ret.Query, newParam(action, target, a, att))
			pv[a] = parseConstraint(kind, att.Validation, att.IsRequired(a))
		}
		if len(pv) > 0 {
			ret.Validator.constraint["payload"] = pv
		}
	}

	// Payload and Query are stored same Query field.
	if action.Payload != nil {
		m := action.Payload.Type.ToObject()
		pv := make(map[string]Constraint)
		for a, att := range m {
			kind := convertTypeString(att.Type.Kind(), target)
			ret.Query = append(ret.Query, newParam(action, target, a, att))
			pv[a] = parseConstraint(kind, att.Validation, att.IsRequired(a))
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
	for _, path := range p.Path {
		c = append(c, fmt.Sprintf("%s(%s): %s", path.Name, path.Kind, path.Description))
	}
	if p.Query != nil {
		c = append(c, fmt.Sprintf("payload(object): payload"))
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
	return p.Name
}

func (p ParamsDefinition) FuncArgs(target string) string {
	ret := make([]string, len(p.Path))
	for i, p := range p.Path {
		switch target {
		case TargetFlow:
			ret[i] = fmt.Sprintf("%s: %s", p.Name, p.Kind)
		case TargetTS:
			ret[i] = fmt.Sprintf("%s: %s", p.Name, p.Kind)
		default:
			ret[i] = p.Name
		}
	}
	if len(p.Query) > 0 {
		var p string
		switch target {
		case TargetFlow:
			p = "payload: any" // TODO
		case TargetTS:
			p = "payload: any" // TODO: should define interface?
		default:
			p = "payload"
		}
		ret = append(ret, p)
	}

	return strings.Join(ret, ", ")
}

func (p ParamsDefinition) UrlArgs() string {
	path := p.Action.Routes[0].FullPath()
	for _, o := range p.Action.Routes[0].Params() {
		d := fmt.Sprintf("${%s}", codegen.Goify(o, false))
		path = strings.Replace(path, ":"+o, d, 1)
	}

	return "`" + path + "`"
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
	Kind        string // kind such as bool, number, ...
	Description string
}

type Params []Param

func newParam(action *design.ActionDefinition, target string, a string, att *design.AttributeDefinition) Param {
	return Param{
		original:    att,
		Name:        codegen.Goify(a, false),
		Kind:        convertTypeString(att.Type.Kind(), target),
		Description: att.Description,
	}
}
