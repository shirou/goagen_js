package goagen_js

import (
	"encoding/json"
	"strings"

	"github.com/goadesign/goa/dslengine"
)

type Validator struct {
	name       string
	original   *dslengine.ValidationDefinition
	constraint map[string]interface{}
}

func (v Validator) JSONify() (string, error) {
	j, err := json.MarshalIndent(v.constraint, " ", "  ")
	if err != nil {
		return "", err
	}
	return string(j), nil
}

type Constraint struct {
	Kind      string        `json:"kind,omitempty"`
	Enum      []interface{} `json:"enum,omitempty"`
	Format    string        `json:"format,omitempty"`
	Pattern   string        `json:"pattern,omitempty"`
	Minimum   *float64      `json:"minimum,omitempty"`
	Maximum   *float64      `json:"maximum,omitempty"`
	MinLength *int          `json:"min_length,omitempty"`
	MaxLength *int          `json:"max_length,omitempty"`
	Required  *bool         `json:"required,omitempty"`
}

func newValidator(name string) Validator {
	return Validator{
		name:       name,
		constraint: make(map[string]interface{}),
	}
}

func parseConstraint(kind string, o *dslengine.ValidationDefinition, required bool) Constraint {
	ret := Constraint{
		Kind: kind,
	}
	// TODO: IsRequired is not work.
	// https://godoc.org/github.com/goadesign/goa/design#AttributeDefinition.IsRequired
	if required {
		ret.Required = &required
	}

	if o == nil {
		return ret
	}

	if o.Values != nil {
		ret.Enum = o.Values
	}
	if o.Format != "" {
		// unescape
		ret.Format = strings.Replace(o.Format, `\`, "", -1)
	}
	if o.Pattern != "" {
		// unescape
		ret.Pattern = strings.Replace(o.Pattern, `\`, "", -1)
	}
	if o.Minimum != nil {
		ret.Minimum = o.Minimum
	}
	if o.Maximum != nil {
		ret.Maximum = o.Maximum
	}
	if o.MinLength != nil {
		ret.MinLength = o.MinLength
	}
	if o.MaxLength != nil {
		ret.MaxLength = o.MaxLength
	}

	return ret
}
