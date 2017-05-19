package goagen_js

import (
	"encoding/json"
	"fmt"

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

func parseConstraint(kind string, o *dslengine.ValidationDefinition) Constraint {
	ret := Constraint{
		Kind: kind,
	}

	if o == nil {
		return ret
	}
	fmt.Printf("%#v\n", o.Required)

	if o.Values != nil {
		ret.Enum = o.Values
	}
	if o.Format != "" {
		ret.Format = o.Format
	}
	if o.Pattern != "" {
		ret.Pattern = o.Pattern
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
	//	v.AddRequired(o.Required)

	return ret
}
