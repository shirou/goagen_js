package goagen_js

const validatorHeaderT = `// This module exports functions that validate {{ .API.Name }} API params hosted at {{ .API.Host }}.
{{ if eq "flow" .Target }}// @flow {{- end }}

export const RequiredError = "missing required parameter";
export const InvalidEnumValueError = "invalid enum value";
export const InvalidFormatError = "invalid format";
export const InvalidPatternError = "invalid pattern";
export const InvalidRangeError = "range exceeded";
export const InvalidMinLengthError = "length is less";
export const InvalidMaxLengthError = "length is exceeded";
export const InvalidKindError = "invalid kind";
`

const validatorT = `
export const {{ .Name }} = {{ .Constraint }};
`

const validatorModuleT = `export function validate(rule, actual) {
  let errors = {};
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
`

const jsFuncsT = `{{ $funcName := .FuncName }}
// {{join .Comments "\n// "}}
export function {{ $funcName }}({{ .Args }}) {
  const url = urlPrefix + {{ .UrlArgs }};

{{- if eq .ValidateRequired true }}
  let e = undefined;
{{- range $p := .PathParams }}
  e = v.validate(v.{{ $funcName }}.{{ $p.Name }}, {{ $p.Name }});
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
`

const jsHeaderT = `// This module exports functions that give access to the {{ .API.Name }} API hosted at {{ .API.Host }}.
{{ if eq "flow" .Target }}// @flow {{- end }}

import 'whatwg-fetch';

import * as v from "./api_validator.js";

const scheme = '{{ .Scheme }}';
const host = '{{ .Host }}';
const urlPrefix = scheme + '://' + host;
`

const jsModuleT = `
// helper function for GET method.
{{      if eq "flow" .Target }}function get(url: string, payload: any): Promise<any> {
{{ else if eq "type" .Target }}function get(url: string, payload: any): Promise<any> {
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
{{      if eq "flow" .Target }}function post(url: string, payload: any): Promise<any> {
{{ else if eq "type" .Target }}function post(url: string, payload: any): Promise<any> {
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
`
