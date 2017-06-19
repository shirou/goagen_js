package goagen_js

import (
	"fmt"
	"strings"
	"text/template"
)

func newTemplate(args ...string) (*template.Template, error) {
	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		"join": func(value []string, arg string) string {
			return strings.Join(value, arg)
		},
	})
	for _, t := range args {
		var err error
		tmpl, err = tmpl.Parse(t)
		if err != nil {
			return nil, fmt.Errorf("tmpl.Parse() failed: %+v, %v", t, err)
		}
	}
	return tmpl, nil
}

const validatorHeaderT = `{{- define "validator_header" -}}
// This module exports functions that validate {{ .API.Name }} API params hosted at {{ .API.Host }}.
{{      if eq "flow" .Target }}// @flow
{{ else if eq "type" .Target }}///<reference path="api.d.ts" />
{{- end }}
{{ end }}
`

const validatorT = `{{- define "validator_definition" }}
export const {{ .Name }} = {{ .Constraint }};
{{ end }}
`

const validatorModuleT = `{{ define "validator_module" }}
export const RequiredError = "missing required parameter";
export const InvalidEnumValueError = "invalid enum value";
export const InvalidFormatError = "invalid format";
export const InvalidPatternError = "invalid pattern";
export const InvalidRangeError = "range exceeded";
export const InvalidMinLengthError = "length is less";
export const InvalidMaxLengthError = "length is exceeded";
export const InvalidKindError = "invalid kind";
{{      if eq "flow" .Target }}export function validate(rule: any, actual: any) {
  let errors = {};
{{ else if eq "type" .Target }}export function validate(rule: any, actual: any): ErrorMap | undefined {
  let errors: ErrorMap = {};
{{ else }}export function validate(rule, actual) {
  let errors = {};
{{ end }}
  if (typeof actual === "object") {
    Object.keys(actual).forEach(function(key, index) {
      const ret = validate(rule[key], actual[key]);
      if (ret !== undefined) {
        errors[key] = ret;
      }
    });
  } else {
    if (rule.kind !== typeof actual){
      errors.kind = InvalidKindError;
      return errors;
    }
    if (rule.maximum && actual > rule.maximum) {
      errors.maximum = InvalidRangeError;
    }
    if (rule.minimum && actual < rule.minimum) {
      errors.minimum = InvalidRangeError;
    }
    if (rule.max_length && actual.length > rule.max_length) {
      errors.max_length = InvalidMaxLengthError;
    }
    if (rule.min_length && actual.length < rule.min_length) {
      errors.min_length = InvalidMinLengthError;
    }
    if (rule.format && !(new RegExp(rule.format).test(actual))) {
      errors.format = InvalidFormatError;
    }
    if (rule.pattern && !(new RegExp(rule.pattern).test(actual))) {
      errors.pattern = InvalidPatternError;
    }
    if (rule.enum) {
      let found = false;
      for (let k in rule.enum) {
        if (k === actual) {found = true;}
      }
      if (found === false) {
        errors.enum = InvalidEnumValueError;
      }
    }
  }
  if (Object.keys(errors).length > 0){
    return errors;
  }
  return undefined;
}
{{- end -}}
`

const jsFuncsT = `{{ define "js_funcs" -}}
{{ $funcName := .FuncName }}
// {{join .Comments "\n// "}}
export function {{ $funcName }}({{ .Args }}){{ .FuncRet }} {
  const url = urlPrefix + {{ .UrlArgs }};

{{- if eq .ValidateRequired true }}
  let e = undefined;
{{- range $p := .PathParams }}
  e = v.validate(v.{{ $funcName }}.{{ $p.CamelCaseName }}, {{ $p.Name }});
  if (e) {
    return Promise.reject(e);
  }
{{- end }}
{{- if .QueryParams }}
  e = v.validate(v.{{ $funcName }}.payload, payload);
  if (e) {
    return Promise.reject(e);
  }
{{- end }}
{{- end }}
  return {{ .Request }};
}
{{- end }}
`

const jsHeaderT = `{{- define "js_header" -}}
// This module exports functions that give access to the {{ .API.Name }} API hosted at {{ .API.Host }}.
{{      if eq "flow" .Target }}// @flow
import * as v from "./api_validator.js";
{{ else if eq "type" .Target }}///<reference path="api.d.ts" />
import * as v from "./api_validator.ts";
{{ else }}
import * as v from "./api_validator.js";
{{- end }}

import 'whatwg-fetch';

const scheme = '{{ .Scheme }}';
const host = '{{ .Host }}';
const urlPrefix = scheme + '://' + host;
{{ end }}
`

const jsModuleT = `{{ define "js_module" }}
// helper function for GET method.
{{      if eq "flow" .Target }}function get(url: string, payload?: any): Promise<any> {
{{ else if eq "type" .Target }}function get(url: string, payload?: any): Promise<any> {
{{ else }}function get(url, payload) {
{{- end }}
  const query = queryBuilder(payload);
  return fetch(url + query, {
    method: 'GET',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
  });
}

// helper function for POST method.
{{      if eq "flow" .Target }}function post(url: string, payload?: any): Promise<any> {
{{ else if eq "type" .Target }}function post(url: string, payload?: any): Promise<any> {
{{ else }}function post(url, payload) {
{{- end }}
  return fetch(url, {
    method: 'POST',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });
}

// helper functon which return QueryParameter from Object.
{{      if eq "flow" .Target }}function queryBuilder(obj: any): string {
{{ else if eq "type" .Target }}function queryBuilder(obj: any): string {
{{ else }}function queryBuilder(obj) {
{{- end }}
  if (!obj) {
    return '';
  }
  const r = Object.keys(obj).sort().map((key) => {
    const val = obj[key];
    if (val === undefined){
      return '';
    }
    if (val === null){
      return '';
    }
    return encodeURIComponent(key) + "=" + encodeURIComponent(val);
  }).filter((x) => {
    return x.length > 0;
  }).join('&');
  if (r.length > 0){
    return '?' + r;
  }
  return '';
}
{{ end -}}
`

const definitionHeaderFlow = `{{ define "definition_header"}}
{{ end }}
`

const definitionFlow = `{{ define "definition"}}
type {{ .Name}}Payload = {
{{- range $p := .PayloadDefinition }}
  {{ $p }}
{{- end }}
};
{{ end }}
`

const definitionHeaderType = `{{ define "definition_header"}}
declare class ErrorMap {
   kind?: string;
   maximum?: string;
   minimum?: string;
   max_length?: string;
   min_length?: string;
   format?: string;
   pattern?: string;
   enum?: string;
   [key: string]: ErrorMap | string | undefined;
}
{{ end }}
`

const definitionType = `{{ define "definition"}}
interface {{ .Name}} {
{{- range $p := .Definition }}
  {{ $p }}
{{- end }}
}
{{ end }}
`
