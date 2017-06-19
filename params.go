package goagen_js

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/goadesign/goa/design"
	"github.com/goadesign/goa/goagen/codegen"
)

const (
	TargetJS   = "js"
	TargetFlow = "flow"
	TargetTS   = "typescript"
)

type ParamsDefinition struct {
	Action    *design.ActionDefinition
	Base      string
	Name      string
	Path      Params // sorted by goa order
	Query     Params // sorted by alphabetical
	Validator Validator
	Response  *Response
}

type Response struct {
	Name           string
	Identifier     string
	IdentifierName string
	Params         Params
}

func parseActions(g *Generator) ([]ParamsDefinition, error) {
	ret := make([]ParamsDefinition, 0)
	actions := getActions(g)

	responses, err := parseMediaTypes(g)
	if err != nil {
		return ret, err
	}

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
			if resp := getResponse(p, responses); resp != nil {
				p.Response = resp
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

func parseMediaTypes(g *Generator) (map[string]Response, error) {
	ret := make(map[string]Response)

	err := g.API.IterateMediaTypes(func(mt *design.MediaTypeDefinition) error {
		if mt.IsError() {
			return nil
		}
		err := mt.IterateViews(func(view *design.ViewDefinition) error {
			params := make(Params, 0)
			for name, att := range view.Type.ToObject() {
				p := Param{
					original:      att,
					Name:          codegen.Goify(name, false),
					CamelCaseName: codegen.Goify(name, true),
					Kind:          convertTypeString(att.Type.Kind(), g.Target),
					Description:   att.Description,
				}
				params = append(params, p)
			}
			ret[mt.Identifier] = Response{
				Identifier:     mt.Identifier,
				IdentifierName: toIdentifierName(mt.Identifier),
				Params:         params,
			}

			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	return ret, err
}

func getResponse(p ParamsDefinition, responses map[string]Response) *Response {
	if p.Action.Responses != nil {
		for _, resp := range p.Action.Responses {
			if resp.MediaType == "" {
				continue
			}
			r, ok := responses[resp.MediaType]
			if !ok {
				continue
			}
			return &r
		}
	}

	return nil
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

func (p ParamsDefinition) PayloadDefinition(target string) []string {
	buf := make([]string, 0)
	for _, tmp := range p.Query {
		switch target {
		case TargetFlow:
			buf = append(buf, fmt.Sprintf("%s: %s,", tmp.Name, tmp.Kind))
		case TargetTS:
			if tmp.Enum == "" {
				buf = append(buf, fmt.Sprintf("%s: %s;", tmp.Name, tmp.Kind))
			} else {
				buf = append(buf, fmt.Sprintf("%s: %s;", tmp.Name, tmp.Enum))
			}
		}
	}

	return buf
}

func (p ParamsDefinition) ResponseDefinition(target string) []string {
	buf := make([]string, 0)
	if p.Response.Params == nil {
		return []string{}
	}
	for _, tmp := range p.Response.Params {
		switch target {
		case TargetFlow:
			buf = append(buf, fmt.Sprintf("%s: %s,", tmp.Name, tmp.Kind))
		case TargetTS:
			if tmp.Enum == "" {
				buf = append(buf, fmt.Sprintf("%s: %s;", tmp.Name, tmp.Kind))
			} else {
				buf = append(buf, fmt.Sprintf("%s: %s;", tmp.Name, tmp.Enum))
			}
		}
	}

	return buf
}

func (p ParamsDefinition) FuncRet(target string) string {
	if p.Response == nil {
		switch target {
		case TargetFlow:
			return ":Promise<any>"
		case TargetTS:
			return ":Promise<any>"
		default:
			return ""
		}
	}
	switch target {
	case TargetFlow:

		return fmt.Sprintf(":Promise<%s>", p.Response.IdentifierName)
	case TargetTS:
		return fmt.Sprintf(":Promise<%s>", p.Response.IdentifierName)
	default:
		return ""
	}

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
		var tmp string
		switch target {
		case TargetFlow:
			tmp = fmt.Sprintf("payload: %sPayload", p.Name)
		case TargetTS:
			tmp = fmt.Sprintf("payload: %sPayload", p.Name)
		default:
			tmp = "payload"
		}
		ret = append(ret, tmp)
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
	original      *design.AttributeDefinition
	Name          string // no CamelCase name
	CamelCaseName string // CamelCase name
	Kind          string // kind such as bool, number, ...
	Description   string
	Enum          string // Enum name
}

type Params []Param

func newParam(action *design.ActionDefinition, target string, a string, att *design.AttributeDefinition) Param {
	p := Param{
		original:      att,
		Name:          codegen.Goify(a, false),
		CamelCaseName: codegen.Goify(a, true),
		Kind:          convertTypeString(att.Type.Kind(), target),
		Description:   att.Description,
	}

	if att.Validation != nil && att.Validation.Values != nil {
		p.Enum = enumValues(att.Validation.Values)
	}

	return p
}

func enumValues(value interface{}) string {
	j, err := json.Marshal(value)
	if err == nil {
		return string(j)
	}
	return "[]"
}

// toIdentifierName convert Identifier to CamelCase Name
// ex: application/vnd.user+json -> UserMedia
// ex: application/vnd.user+json; type=collection -> CollectionOfUserMedia
func toIdentifierName(identifier string) string {
	tmp := design.CanonicalIdentifier(identifier)
	ret := strings.Replace(tmp, "application/vnd.", "", 1)
	ret = codegen.Goify(ret, true) + "Media"

	if strings.Contains("type=collection", tmp) {
		ret = "CollectionOf" + ret
	}

	return ret
}
