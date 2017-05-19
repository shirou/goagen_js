// This module exports functions that validate goagen_js API params hosted at localhost:8080.

export const RequiredError = "missing required parameter";
export const InvalidEnumValueError = "invalid enum value";
export const InvalidFormatError = "invalid format";
export const InvalidPatternError = "invalid pattern";
export const InvalidRangeError = "range exceeded";
export const InvalidLengthError = "length is exceeded or less";
export const InvalidKindError = "invalid kind";

export const GetIntGet = {
   "payload": {
     "int": {
       "kind": "number"
     },
     "int_array": {
       "kind": "array"
     },
     "int_enum": {
       "kind": "number",
       "enum": [
         1,
         2,
         3
       ]
     },
     "int_max": {
       "kind": "number",
       "maximum": 10
     },
     "int_min": {
       "kind": "number",
       "minimum": -1
     },
     "int_minmax": {
       "kind": "number",
       "minimum": 0,
       "maximum": 10
     },
     "int_required": {
       "kind": "number"
     },
     "int_secret": {
       "kind": "number"
     }
   }
 };

export const PathParamsGet = {
   "ParamInt": {
     "kind": "number",
     "maximum": 10
   },
   "ParamStr": {
     "kind": "string"
   }
 };

export const WithoutGet = {};
export function validate(rule, actual) {
  let errors = {};

  if (typeof actual === "object") {
    Object.keys(actual).map(function(key, index) {
      errors[key] = validate(rule[key], actual[key]);
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

  if (errors.length === 0) {
    return null;
  }
  return errors;
}
