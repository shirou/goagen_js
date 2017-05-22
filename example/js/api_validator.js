// This module exports functions that validate goagen_js API params hosted at localhost:8080.

export const RequiredError = "missing required parameter";
export const InvalidEnumValueError = "invalid enum value";
export const InvalidFormatError = "invalid format";
export const InvalidPatternError = "invalid pattern";
export const InvalidRangeError = "range exceeded";
export const InvalidLengthError = "length is exceeded or less";
export const InvalidKindError = "invalid kind";

export const UserCreate = {};

export const UserGet = {
   "UserID": {
     "kind": "number",
     "maximum": 10000
   }
 };

export const UserList = {};
export function validate(rule, actual) {
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
