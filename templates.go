package goagen_js

const validatorHeaderT = `// This module exports functions that validate {{ .API.Name }} API params hosted at {{ .API.Host }}.

export const RequiredError = "missing required parameter";
export const InvalidEnumValueError = "invalid enum value";
export const InvalidFormatError = "invalid format";
export const InvalidPatternError = "invalid pattern";
export const InvalidRangeError = "range exceeded";
export const InvalidLengthError = "length is exceeded or less";
export const InvalidKindError = "invalid kind";
`

const validatorT = `
export const {{ .Name }} = {{ .Constraint }};
`

const validatorModuleT = `export function validate(rule, actual) {
  let errors = {};

  if (typeof actual === "object") {
    Object.keys(actual).map(function(key, index) {
      const ret = validate(rule[key], actual[key]);
      if (ret !== null) {
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
      errors.maximum = InvalidRangeError;
    }
    if (rule.max_length && actual.length < rule.max_length) {
      errors.max_length = InvalidLengthError;
    }
    if (rule.min_length && actual.length < rule.min_length) {
      errors.min_length = InvalidLengthError;
    }
    if (rule.format && new RegExp(rule.format).test(actual)) {
      errors.format = InvalidFormatError;
    }
    if (rule.pattern && new RegExp(rule.pattern).test(actual)) {
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
  if (Object.keys(errors).length === 0) {
    return null;
  }
  return errors;
}
`

const jsFuncsT = `{{ $funcName := .FuncName }}
// {{join .Comments "\n// "}}
export function {{ $funcName }}({{ .Args }}) {
  const url = urlPrefix + {{ .UrlArgs }};

{{- if eq .ValidateRequired true }}
  let errors = {};
  let ret;
{{- range $p := .PathParams }}
  if (v.validate(v.{{ $funcName }}.{{ $p.Name }}, {{ $p.Name }}) !== null) {
    return Promise.reject(new Error("validation error"));
  }
{{- end }}
{{- if .QueryParams }}
  if (v.validate(v.{{ $funcName }}.payload, payload) !== null) {
    return Promise.reject(new Error("validation error"));
  }
{{- end }}
{{- end }}
  return {{ .Request }};
}
`

const jsHeaderT = `// This module exports functions that give access to the {{ .API.Name }} API hosted at {{ .API.Host }}.

import 'whatwg-fetch';

import * as v from "./api_validator.js";

const scheme = '{{ .Scheme }}';
const host = '{{ .Host }}';
const urlPrefix = scheme + '://' + host;
`

const jsModuleT = `
// helper function for GET method.
function get(url, payload) {
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
function post(url, payload) {
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
function queryBuilder(obj) {
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
